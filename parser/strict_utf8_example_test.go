// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser_test

import (
	"errors"
	"fmt"

	"github.com/leodido/go-conventionalcommits/parser"
)

func ExampleMachine_WithStrictUTF8() {
	machine := parser.NewMachine()
	machine.WithStrictUTF8()

	message, err := machine.Parse(append([]byte("feat: invalid "), 0xff))

	var invalidUTF8Error *parser.InvalidUTF8Error
	fmt.Println(message == nil)
	fmt.Println(errors.Is(err, parser.ErrInvalidUTF8))
	if errors.As(err, &invalidUTF8Error) {
		fmt.Println(invalidUTF8Error.ByteOffset)
	}
	// Output:
	// true
	// true
	// 14
}
