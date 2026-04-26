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

// WithStrictUTF8 enables opt-in UTF-8 well-formedness validation of
// trailer values, scope, and free-form types.
//
// When enabled, after a successful parse the captured slices are
// checked with utf8.Valid. The first ill-formed slice produces an
// error whose column is the byte offset (in the ORIGINAL input) of
// the first invalid byte; trailer-value errors also include the
// failing footer key.
//
// Body and description are intentionally NOT covered by this option.
// Use WithStrictUTF8Body for those.
func WithStrictUTF8() conventionalcommits.MachineOption {
	return conventionalcommits.WithStrictUTF8()
}

// WithStrictUTF8Body enables opt-in UTF-8 well-formedness validation
// of body and description slices.
//
// When enabled, after a successful parse the captured slices are
// checked with utf8.Valid. The first ill-formed slice produces an
// error whose column is the byte offset (in the ORIGINAL input) of
// the first invalid byte.
//
// Trailer values, scope, and free-form types are intentionally NOT
// covered by this option. Use WithStrictUTF8 for those.
func WithStrictUTF8Body() conventionalcommits.MachineOption {
	return conventionalcommits.WithStrictUTF8Body()
}
