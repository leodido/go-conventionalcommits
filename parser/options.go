// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"github.com/leodido/go-conventionalcommits"
	"github.com/sirupsen/logrus"
)

// Machine is a Conventional Commits parser with parser-specific configuration.
type Machine interface {
	conventionalcommits.Machine

	// WithStrictUTF8 requires the entire original input to be well-formed UTF-8.
	// Malformed input returns a nil message and an error at the first invalid byte,
	// including when best effort mode is enabled. The error matches ErrInvalidUTF8
	// and can be inspected as *InvalidUTF8Error.
	WithStrictUTF8()
}

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

func (m *machine) WithStrictUTF8() {
	m.strictUTF8 = true
}
