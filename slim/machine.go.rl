package slim

import (
	"fmt"

	"github.com/leodido/go-conventionalcommits"
)

// ColumnPositionTemplate is the template used to communicate the column where errors occur.
var ColumnPositionTemplate = ": col=%02d"

const (
	// ErrType represents an error in the type part of the commit message.
	ErrType = "illegal '%s' character in commit message type"
	// ErrColon is the error message that communicate that the mandatory colon after the type part of the commit message is missing.
	ErrColon = "expecting colon (':') character, got '%s' character"
	// ErrTypeIncomplete represents an error in the type part of the commit message.
	ErrTypeIncomplete = "incomplete commit message type after '%s' character"
	// ErrMalformedScope represents an error about illegal characters into the the scope part of the commit message.
	ErrMalformedScope = "illegal '%s' character in scope"
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
	// ErrEarly represents an error when the input makes the machine exit too early.
	ErrEarly = "early exit after '%s' character"
	// ErrDescriptionInit tells the user that before of the description part a whitespace is mandatory.
	ErrDescriptionInit = "expecting at least one white-space (' ') character, got '%s' character"
	// ErrDescription tells the user that after the whitespace is mandatory a description.
	ErrDescription = "expecting a description after '%s' character"
)

%%{
machine conventionalcommits;

include common "common.rl";

# unsigned alphabet
alphtype uint8;

action mark {
	m.pb = m.p
}

# Error management

action err_empty {
	m.err = m.emitErrorWithoutCharacter(ErrEmpty)
}

action err_type {
	if m.pe > 0 {
		if m.p != m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrType)
		} else {
			m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
		}
	}
}

action err_malformed_scope {
	m.err = m.emitErrorOnCurrentCharacter(ErrMalformedScope)
}

action err_colon {
	if m.err == nil {
		m.err = m.emitErrorOnCurrentCharacter(ErrColon);
	}
}

action err_description_init {
	if m.err == nil {
		m.err = m.emitErrorOnCurrentCharacter(ErrDescriptionInit);
	}
}

action err_description {
	m.err = m.emitErrorOnPreviousCharacter(ErrDescription);
}

action check_early_exit {
	if (m.p + 1) == m.pe {
		m.err = m.emitErrorOnCurrentCharacter(ErrEarly);
	}
}

# Setters

action set_type {
	output._type = string(m.text())
}

action set_scope {
	output.scope = string(m.text())
}

action set_description {
	output.descr = string(m.text())
}

action set_exclamation {
	output.exclamation = true
}

# Machine definitions

minimal_types = ('fix' | 'feat');

conventional_types = ('build' | 'ci' | 'chore' | 'docs' | 'feat' | 'fix' | 'perf' | 'refactor' | 'revert' | 'style' | 'test');

falco_types = ('build' | 'ci' | 'chore' | 'docs' | 'feat' | 'fix' | 'perf' | 'new' | 'revert' | 'update' | 'test' | 'rule' );

scope = lpar ((any* -- lpar) -- rpar) >mark %err(err_malformed_scope) %set_scope rpar;

breaking = exclamation >set_exclamation;

## todo > strict option to enforce a single whitespace?
description = ws+ >err(err_description_init) <: any+ >mark >err(err_description) %set_description;

## todo > option to limit the total length
main := minimal_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	scope? %to(check_early_exit)
	breaking? %to(check_early_exit)
	colon >err(err_colon) %to(check_early_exit)
	description;

conventional_types_main := conventional_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	scope? %to(check_early_exit)
	breaking? %to(check_early_exit)
	colon >err(err_colon) %to(check_early_exit)
	description;

falco_types_main := falco_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	scope? %to(check_early_exit)
	breaking? %to(check_early_exit)
	colon >err(err_colon) %to(check_early_exit)
	description;

}%%

%% write data noerror noprefix;

type machine struct {
	data         []byte
	cs           int
	p, pe, eof   int
	pb           int
	err          error
	bestEffort   bool
	typeConfig   conventionalcommits.TypeConfig
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

func (m *machine) emitErrorWithoutCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate + ColumnPositionTemplate, m.p)
}

func (m *machine) emitErrorOnCurrentCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate + ColumnPositionTemplate, string(m.data[m.p]), m.p)
}

func (m *machine) emitErrorOnPreviousCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate + ColumnPositionTemplate, string(m.data[m.p - 1]), m.p)
}

// NewMachine creates a new FSM able to parse Conventional Commits.
func NewMachine(options ...conventionalcommits.MachineOption) conventionalcommits.Machine {
	m := &machine{}

	for _, opt := range options {
		opt(m)
	}

	%% access m.;
	%% variable p m.p;
	%% variable pe m.pe;
	%% variable eof m.eof;
	%% variable data m.data;

	return m
}

// Parse parses the input byte array as a Conventional Commit message with no body neither footer.
//
// When a valid Conventional Commit message is given it outputs its structured representation.
// If the parsing detects an error it returns it with the position where the error occurred.
//
// It can also partially parse input messages returning a partially valid structured representation
// and the error that stopped the parsing.
func (m *machine) Parse(input []byte) (conventionalcommits.Message, error) {
	m.data = input
	m.p = 0
	m.pb = 0
	m.pe = len(input)
	m.eof = len(input)
	m.err = nil
	output := &conventionalCommit{}

	switch m.typeConfig {
	case conventionalcommits.TypesConventional:
		m.cs = en_conventional_types_main
		break
	case conventionalcommits.TypesFalco:
		m.cs = en_falco_types_main
		break
	case conventionalcommits.TypesMinimal:
		fallthrough
	default:
		%% write init;
		break
	}
	%% write exec;

	if m.cs < first_final {
		if m.bestEffort && output.minimal() {
			// An error occurred but partial parsing is on and partial message is minimally valid
			return output.export(), m.err
		}
		return nil, m.err
	}

	return output.export(), nil
}

// WithBestEffort enables best effort mode.
func (m *machine) WithBestEffort() {
	m.bestEffort = true
}

// HasBestEffort tells whether the receiving machine has best effort mode on or off.
func (m *machine) HasBestEffort() bool {
	return m.bestEffort
}

// WithTypes ...
func (m *machine) WithTypes(t conventionalcommits.TypeConfig) {
	m.typeConfig = t
}