package slim

import (
	"github.com/leodido/go-conventionalcommits"
)

// WithBestEffort enables the best effort mode.
func WithBestEffort() conventionalcommits.MachineOption {
	return func(m conventionalcommits.Machine) conventionalcommits.Machine {
		m.WithBestEffort()
		return m
	}
}
