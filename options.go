// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package conventionalcommits

import (
	"github.com/sirupsen/logrus"
)

// WithBestEffort ...
func WithBestEffort() MachineOption {
	return func(m Machine) Machine {
		m.(BestEfforter).WithBestEffort()

		return m
	}
}

// WithTypes ...
func WithTypes(t TypeConfig) MachineOption {
	return func(m Machine) Machine {
		m.(TypeConfigurer).WithTypes(t)

		return m
	}
}

// WithLogger ...
func WithLogger(l *logrus.Logger) MachineOption {
	return func(m Machine) Machine {
		m.(Logger).WithLogger(l)

		return m
	}
}

// WithStrictUTF8 enables opt-in UTF-8 well-formedness validation of
// the slices captured for trailer values, scope, and free-form
// types.
//
// The check is a no-op for Machine implementations that do not
// implement StrictUTF8er (no panic, no error).
//
// Body and description are intentionally NOT covered by this option.
// Use WithStrictUTF8Body for those (the two knobs are independent
// and can be combined).
func WithStrictUTF8() MachineOption {
	return func(m Machine) Machine {
		if x, ok := m.(StrictUTF8er); ok {
			x.WithStrictUTF8()
		}

		return m
	}
}

// WithStrictUTF8Body enables opt-in UTF-8 well-formedness validation
// of body and description slices.
//
// The check is a no-op for Machine implementations that do not
// implement StrictUTF8Bodyer (no panic, no error).
//
// Trailer values, scope, and free-form types are intentionally NOT
// covered by this option. Use WithStrictUTF8 for those.
func WithStrictUTF8Body() MachineOption {
	return func(m Machine) Machine {
		if x, ok := m.(StrictUTF8Bodyer); ok {
			x.WithStrictUTF8Body()
		}

		return m
	}
}
