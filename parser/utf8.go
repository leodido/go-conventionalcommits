// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import "unicode/utf8"

const (
	// ErrInvalidUTF8 represents malformed UTF-8 anywhere in the original input.
	ErrInvalidUTF8 = "invalid UTF-8"
)

// firstInvalidUTF8Index returns the byte index of the first malformed UTF-8
// byte, or -1 when input is well-formed. A correctly encoded U+FFFD is valid:
// DecodeRune reports it with a width greater than one.
func firstInvalidUTF8Index(input []byte) int {
	for i := 0; i < len(input); {
		r, size := utf8.DecodeRune(input[i:])
		if r == utf8.RuneError && size == 1 {
			return i
		}
		i += size
	}

	return -1
}
