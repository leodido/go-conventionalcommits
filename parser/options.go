// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"github.com/leodido/go-conventionalcommits"
	"github.com/sirupsen/logrus"
)

// WithBestEffort enables the best effort mode.
//
// Best effort mode tells the parser to return what it found,
// if the input was a minimally well-formed commit message (type and description part).
func WithBestEffort() conventionalcommits.MachineOption {
	return func(m conventionalcommits.Machine) conventionalcommits.Machine {
		m.WithBestEffort()

		return m
	}
}

// WithTypes let you choose the types.
func WithTypes(t conventionalcommits.TypeConfig) conventionalcommits.MachineOption {
	return func(m conventionalcommits.Machine) conventionalcommits.Machine {
		m.WithTypes(t)

		return m
	}
}

// WithLogger enables a logger during parsing.
func WithLogger(l *logrus.Logger) conventionalcommits.MachineOption {
	return func(m conventionalcommits.Machine) conventionalcommits.Machine {
		m.WithLogger(l)

		return m
	}
}

// WithStrictUTF8 requires the entire commit message to be well-formed UTF-8.
//
// Validation runs on the original input before parsing. Malformed input returns
// a nil message and an error at the first invalid byte, including when best
// effort mode is enabled.
//
// This option is specific to machines created by this parser package. Applying
// it to another Machine implementation panics instead of silently doing nothing.
func WithStrictUTF8() conventionalcommits.MachineOption {
	return func(m conventionalcommits.Machine) conventionalcommits.Machine {
		parserMachine, ok := m.(*machine)
		if !ok || parserMachine == nil {
			panic("parser.WithStrictUTF8 requires parser.NewMachine")
		}
		parserMachine.strictUTF8 = true

		return parserMachine
	}
}
