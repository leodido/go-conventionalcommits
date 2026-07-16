// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package conventionalcommits

import (
	"github.com/sirupsen/logrus"
)

// TypeConfig represent the set of types the parser should use.
type TypeConfig int

const (
	// TypesMinimal is the minimal set of types.
	TypesMinimal TypeConfig = iota
	// TypesConventional represents the conventional set of types.
	// See https://github.com/conventional-changelog/commitlint/tree/master/%40commitlint/config-conventional
	TypesConventional
	// TypesFalco represents the set of types that Falco uses for its release notes.
	// See https://github.com/falcosecurity/falco
	TypesFalco
	// TypesFreeForm represents a free-form set of types.
	TypesFreeForm
)

// TypeConfigurer represents parsers with the option to enable different commit message types.
type TypeConfigurer interface {
	WithTypes(t TypeConfig)
}

// VersionBump represent the set of possible version bumps a commit can mandate.
type VersionBump int

const (
	UnknownVersion VersionBump = iota
	PatchVersion
	MinorVersion
	MajorVersion
)

// BestEfforter is an interface that wraps the methods about the best effort mode.
type BestEfforter interface {
	WithBestEffort()
	HasBestEffort() bool
}

// Logger represents parser able to log.
type Logger interface {
	WithLogger(l *logrus.Logger)
}

// Machine represent a FSM able to parse a conventional commit and return it in an structured way.
type Machine interface {
	Parse(input []byte) (*ConventionalCommit, error)
	BestEfforter
	TypeConfigurer
	Logger
}

// MachineOption represents the type of option setters for Machine instances.
type MachineOption func(m Machine) Machine

// Message represent a conventional commit message.
type Message interface {
	Ok() bool
	IsBreakingChange() bool
	IsFeat() bool
	IsFix() bool
	VersionBump(VersionBumpStrategy) VersionBump
	HasFooter() bool
}

// ConventionalCommit represents a commit message as per Conventional Commits specification.
type ConventionalCommit struct {
	Type        string
	Description string
	Scope       *string // optional
	Exclamation bool
	Body        *string             // optional
	Footers     map[string][]string // optional
	TypeConfig  TypeConfig
}

// VersionBumpStrategy represents a strategy how to evaluate the version bump depending on the TypeConfig initially used and the commits type.
type VersionBumpStrategy func(*ConventionalCommit) VersionBump

// DefaultStrategy is a basic, opiniated strategy to evaluate the version bump.
func DefaultStrategy(c *ConventionalCommit) VersionBump {
	if c.IsBreakingChange() {
		return MajorVersion
	}
	if c.IsFeat() {
		return MinorVersion
	}
	if c.IsFix() {
		return PatchVersion
	}

	return UnknownVersion
}

// Ok tells whether the receiving commit message is well-formed or not.
//
// A minimally well-formed commit message has at least a valid type and a non empty description.
func (c *ConventionalCommit) Ok() bool {
	return c.Type != "" && c.Description != ""
}

// IsBreakingChange tells whether the receiving commit message struct represents a breaking change or not.
func (c *ConventionalCommit) IsBreakingChange() bool {
	_, hasBreakingChangeTrailer := c.Footers["breaking-change"]

	return c.Exclamation || hasBreakingChangeTrailer
}

// IsFeat tells whether the receiving commit message struct represents a feat change or not.
func (c *ConventionalCommit) IsFeat() bool {
	if c.TypeConfig == TypesFalco && c.Type == "new" {
		return true
	}

	return c.Type == "feat"
}

// IsFix tells whether the receiving commit message struct represents a fix change or not.
func (c *ConventionalCommit) IsFix() bool {
	return c.Type == "fix"
}

// VersionBump tells which version bump the receiving commit message mandates.
func (c *ConventionalCommit) VersionBump(strategy VersionBumpStrategy) VersionBump {
	if strategy == nil {
		return DefaultStrategy(c)
	}

	return strategy(c)
}

// HasFooter tells whether the receiving commit message struct has one or more trailers.
func (c *ConventionalCommit) HasFooter() bool {
	return len(c.Footers) > 0
}
