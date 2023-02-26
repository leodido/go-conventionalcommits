// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package testing

import (
	"strings"
)

// StringAddress returns the address of the input string.
func StringAddress(str string) *string {
	return &str
}

// RightPad pads a string with spaces until the given limit, or it cuts the string to the given limit.
func RightPad(str string, limit int) string {
	str += strings.Repeat(" ", limit)

	return str[:limit]
}
