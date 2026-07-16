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
