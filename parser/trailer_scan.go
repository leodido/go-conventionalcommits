// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"bytes"
)

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
// Performance: the implementation is one forward pass over data to
// build the per-line offset table, plus one reverse pass over that
// table that does at most one isTrailerStartLine check per blank-line
// gap. Both passes are O(n) where n = len(data); no quadratic
// LastIndex behavior, even on inputs with many trailer-shaped
// paragraphs.
func findTrailerBlockStart(data []byte) int {
	if len(data) == 0 {
		return -1
	}

	// Drop trailing newlines so they don't get mis-attributed.
	end := len(data)
	for end > 0 && data[end-1] == '\n' {
		end--
	}
	if end == 0 {
		return -1
	}

	// Forward pass: record the byte offset of the first byte of every
	// non-blank line in data[:end] and whether that line is preceded
	// by a blank-line gap. The first line (i == 0) is the description,
	// blankBefore[0] is always false by convention.
	lineStart, blankBefore := scanLines(data, end)
	if len(lineStart) == 0 {
		return -1
	}

	// Reverse pass: walk lines from the bottom up. Each blank-line gap
	// either extends the trailer block (when the line right after the
	// gap, i.e. the line at lineStart[i], is trailer-shaped) or
	// terminates it. Continuation lines (no blank gap above them) are
	// part of whatever block the line above them belongs to, so we
	// don't classify them.
	candidate := -1
	for i := len(lineStart) - 1; i >= 0; i-- {
		if !blankBefore[i] {
			continue
		}
		if i == 0 {
			// Would extend up to the description line; never classify
			// the description as part of the trailer block.
			return candidate
		}
		if !isTrailerStartLine(lineAt(data, lineStart[i], end)) {
			// The line at lineStart[i] is not trailer-shaped. The
			// previous candidate (set at a deeper line) is the answer.
			return candidate
		}
		// Trailer-shaped first line of a sub-block: extend the block
		// up to here and look for an even deeper gap.
		candidate = lineStart[i]
	}

	return candidate
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
