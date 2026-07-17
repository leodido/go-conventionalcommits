// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2026- Leonardo Di Donato <leodidonato@gmail.com>

package parser

// trimTerminalLineEndings removes a trailing sequence of LF and CRLF tokens
// after non-whitespace content. It deliberately leaves trailing whitespace
// and non-empty lines for the parser to reject normally.
func trimTerminalLineEndings(input []byte) []byte {
	end := len(input)
	for end > 0 && input[end-1] == '\n' {
		end--
		if end > 0 && input[end-1] == '\r' {
			end--
		}
	}
	if end == len(input) || end == 0 || input[end-1] == ' ' || input[end-1] == '\t' {
		return input
	}

	return input[:end]
}
