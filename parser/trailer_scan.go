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

	// Walk bottom-up across blank-line gaps. A gap may extend the
	// trailer block iff the line above the gap is itself trailer-
	// shaped. Otherwise the gap is the body/trailer separator and the
	// block (if any) starts at the first trailer-shaped line below the
	// gap.
	//
	// To enter the loop with a candidate block, the first trailer-
	// shaped line above the bottom-most blank-line gap must exist and
	// must NOT be the description line.
	candidate := -1
	gapEnd := end
	for {
		// Find the previous `\n\n` (blank-line separator) inside
		// data[:gapEnd]. The line that starts immediately after the
		// `\n\n` is the candidate first line of a (possibly extended)
		// trailer block.
		idx := bytes.LastIndex(data[:gapEnd], []byte("\n\n"))
		if idx < 0 {
			// No blank-line separator above; whatever we have so far
			// is the answer.
			return candidate
		}
		// Skip extra newlines: a run of three or more newlines is one
		// gap. The first line after the run starts at `lineStart`.
		lineStart := idx + 2
		for lineStart < end && data[lineStart] == '\n' {
			lineStart++
		}
		if lineStart >= end {
			// All newlines up to EOF; nothing classifiable.
			return candidate
		}
		if lineStart == 0 {
			// Would start at the description line.
			return candidate
		}
		if !isTrailerStartLine(lineAt(data, lineStart, end)) {
			// First line after the gap isn't trailer-shaped: this gap
			// terminates the block at the previous candidate.
			return candidate
		}
		// This gap is inside a trailer block: extend the block up to
		// `lineStart` and look for an even earlier gap.
		candidate = lineStart
		gapEnd = idx
	}
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
