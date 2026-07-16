// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

// ErrInvalidUTF8 identifies malformed UTF-8 anywhere in the original input.
var ErrInvalidUTF8 = errors.New("invalid UTF-8")

// InvalidUTF8Error reports the first malformed byte in the original input.
type InvalidUTF8Error struct {
	// ByteOffset is the zero-based offset of the first malformed byte.
	ByteOffset int
}

func (e *InvalidUTF8Error) Error() string {
	return fmt.Sprintf("%s"+ColumnPositionTemplate, ErrInvalidUTF8, e.ByteOffset)
}

func (e *InvalidUTF8Error) Unwrap() error {
	return ErrInvalidUTF8
}

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
