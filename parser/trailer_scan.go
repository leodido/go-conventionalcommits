// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"bytes"
)

// findTrailerBlock returns both the start offset of the trailing
// trailer block (or -1 if absent) AND the sorted list of byte offsets
// of every line within that block whose first line is itself
// trailer_init-shaped (the "real trailer lines"). The line-starts
// table is empty when there is no trailer block.
//
// The line-starts table is what the FSM consults at runtime via
// trailerValueContinues to enforce spec clause-10 termination of
// multi-line trailer values: a value continues across a newline iff
// the line right after the newline is NOT registered here. See issue
// #48.
//
// findTrailerBlock is the canonical entry point for the parser; the
// older findTrailerBlockStart is preserved as a thin wrapper for the
// existing trailer_scan unit tests.
func findTrailerBlock(data []byte) (start int, lineStarts []int) {
	if len(data) == 0 {
		return -1, nil
	}

	// Drop trailing newlines so they don't get mis-attributed.
	end := len(data)
	for end > 0 && data[end-1] == '\n' {
		end--
	}
	if end == 0 {
		return -1, nil
	}

	// Forward pass: record offsets of every non-blank line's first
	// byte and whether each line is preceded by a blank-line gap.
	allLineStart, blankBefore := scanLines(data, end)
	if len(allLineStart) == 0 {
		return -1, nil
	}

	// Reverse pass: walk from the bottom, extending the candidate
	// trailer block across blank-line gaps when the line above the
	// gap is itself trailer-shaped. This locates the block start.
	candidate := -1
	candidateLineIdx := -1
	for i := len(allLineStart) - 1; i >= 0; i-- {
		if !blankBefore[i] {
			continue
		}
		if i == 0 {
			break
		}
		if !isTrailerStartLine(lineAt(data, allLineStart[i], end)) {
			break
		}
		candidate = allLineStart[i]
		candidateLineIdx = i
	}
	if candidate < 0 {
		return -1, nil
	}

	// Forward pass over just the lines inside the trailer block:
	// classify each as a real-trailer line or a continuation line of
	// the previous trailer's value. Only real-trailer lines go into
	// the lineStarts table.
	lineStarts = make([]int, 0, len(allLineStart)-candidateLineIdx)
	for i := candidateLineIdx; i < len(allLineStart); i++ {
		if isTrailerStartLine(lineAt(data, allLineStart[i], end)) {
			lineStarts = append(lineStarts, allLineStart[i])
		}
	}

	return candidate, lineStarts
}

// findTrailerBlockStart returns the byte offset, within data, at which
// the trailing footer trailer block begins, or -1 if no such block is
// present.
//
// The "trailer block" is the longest run of trailer-shaped sections at
// the end of the message. A "section" is one trailer-shaped first line
// followed by zero or more continuation lines. Sections inside the
// block may be separated from each other by one or more blank lines —
// this matches the FSM's `trailer_beg = nl* trailer_init` semantics,
// where arbitrary blank-line gaps between trailer sub-blocks are
// tolerated. The block ends (going up) at the first blank-line gap
// whose line above is NOT trailer-shaped: that line above is body
// content, and the blank line is the body/trailer separator.
//
// Trailer-shaped, per `trailer_init` in machine.go.rl, means:
//
//   - `alnum+ (- alnum+)*` then either `: <ws>+` (single-space `ws`,
//     per `common.rl`) or ` #<value>`; or
//   - the special-case `BREAKING CHANGE: <ws>+`.
//
// Returning a nonzero offset does NOT assert that every line in the
// block is well-formed: the FSM is responsible for rejecting malformed
// trailer values via its existing rewind/error paths once it enters
// trailer parsing.
//
// The first line of the message (offset 0) is the description and is
// never part of a trailer block, even when its text looks trailer-
// shaped.
//
// Trailing newlines/whitespace at the end of the input are ignored
// when locating the block.
//
// Performance: one forward pass over data to build the per-line
// offset table, plus one reverse pass over that table that does at
// most one isTrailerStartLine check per blank-line gap, plus one final
// forward pass over the in-block lines. All three are O(n) where n =
// len(data); no quadratic LastIndex behavior, even on inputs with
// many trailer-shaped paragraphs.
func findTrailerBlockStart(data []byte) int {
	start, _ := findTrailerBlock(data)

	return start
}

// scanLines walks data[:end] once and returns:
//   - lineStart: the byte offset (within data) of the first byte of
//     each non-blank line, in order;
//   - blankBefore: a parallel slice where blankBefore[i] is true iff
//     the line starting at lineStart[i] is preceded by at least one
//     blank-line gap (i.e., a run of two or more consecutive newlines).
//
// blankBefore[0] is always false (the description has nothing before
// it).
func scanLines(data []byte, end int) ([]int, []bool) {
	// A typical Conventional Commit has between 1 and ~30 lines; pre-
	// allocating a small backing array avoids the first few growths
	// without overshooting on small inputs.
	lineStart := make([]int, 0, 16)
	blankBefore := make([]bool, 0, 16)

	pos := 0
	consecutiveNL := 0
	atLineStart := true
	for pos < end {
		b := data[pos]
		if b == '\n' {
			consecutiveNL++
			atLineStart = true
			pos++

			continue
		}
		if atLineStart {
			lineStart = append(lineStart, pos)
			blankBefore = append(blankBefore, pos != 0 && consecutiveNL >= 2)
			atLineStart = false
			consecutiveNL = 0
		}
		pos++
	}

	return lineStart, blankBefore
}

// lineAt returns the slice of data starting at start and ending at the
// next newline (exclusive) or at end (exclusive), whichever comes first.
func lineAt(data []byte, start, end int) []byte {
	stop := bytes.IndexByte(data[start:end], '\n')
	if stop < 0 {
		return data[start:end]
	}

	return data[start : start+stop]
}

// isTrailerStartLine reports whether line begins a trailer (token plus
// separator), matching `trailer_init` in machine.go.rl. The whitespace
// alphabet is a single ASCII space, matching `ws = ' '` in common.rl.
//
// Grammar:
//
//	trailer_tok = alnum+ (- alnum+)*
//	trailer_sep = (':' ws+) | (ws '#')
//	trailer_init = ('BREAKING CHANGE' >mark trailer_sep_breaking)
//	             | (trailer_tok >mark trailer_sep)
//	trailer_sep_breaking = ':' ws+
//
// Returning true does NOT mean the rest of the line is a valid trailer
// value; the FSM catches malformed values at parse time. The pre-scan
// only needs to confirm the line OPENS like a trailer.
func isTrailerStartLine(line []byte) bool {
	if bytes.HasPrefix(line, []byte("BREAKING CHANGE")) {
		rest := line[len("BREAKING CHANGE"):]
		// trailer_sep_breaking := colon ws+
		if len(rest) >= 2 && rest[0] == ':' && rest[1] == ' ' {
			return true
		}

		return false
	}

	// trailer_tok := alnum+ (dash alnum+)*
	i := 0
	if i >= len(line) || !isAlnum(line[i]) {
		return false
	}
	for i < len(line) && isAlnum(line[i]) {
		i++
	}
	for i < len(line) && line[i] == '-' {
		i++
		if i >= len(line) || !isAlnum(line[i]) {
			return false
		}
		for i < len(line) && isAlnum(line[i]) {
			i++
		}
	}
	if i >= len(line) {
		return false
	}
	// trailer_sep := (colon ws+) | (ws '#')
	switch line[i] {
	case ':':
		i++
		if i >= len(line) || line[i] != ' ' {
			return false
		}

		return true
	case ' ':
		i++
		if i >= len(line) || line[i] != '#' {
			return false
		}

		return true
	}

	return false
}

// isAlnum mirrors Ragel's ASCII-only `alnum` class (digits + ASCII
// letters). If the FSM grammar ever broadens to accept Unicode letters
// this needs to track that change to keep the pre-scan and the FSM in
// agreement.
func isAlnum(b byte) bool {
	return (b >= '0' && b <= '9') || (b >= 'A' && b <= 'Z') || (b >= 'a' && b <= 'z')
}
