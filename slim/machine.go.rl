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
	// ErrTypeIncomplete represents an error in the type part of the commit message.
	ErrTypeIncomplete = "incomplete commit message type after '%s' character"
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
)

%%{
machine conventionalcommits;

include common "common.rl";

# unsigned alphabet
alphtype uint8;

action err_empty {
	m.err = fmt.Errorf(ErrEmpty+ColumnPositionTemplate, m.p)
	fhold;
}

action mark {
	m.pb = m.p
}

action set_type {
	output._type = string(m.text())
}

action err_type {
	if m.p != m.pe {
		m.err = m.emitErrorOnCurrentCharacter(ErrType)
	} else {
		m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
	}
	fhold;
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

## todo > how to configure these? ideally, at runtime...
## todo > error management
type = ('fix' | 'feat') >mark %set_type <err(err_type) >eof(err_empty);

## todo > option to exclude whitespaces and parentheses from valid scope corpus
## todo > error management
scope = op any* >mark %set_scope cp;

## todo > error management
## todo > set description
desc = any* >mark %set_description;

# a machine that consumes the rest of the line when parsing fails
fail := (any - [\n\r])*;

## todo > option to limit the total length
## todo > err_generic is likely unneded
main := type scope? (exclamation >mark %set_exclamation)? colon ws desc;

}%%

%% write data noerror noprefix;

type machine struct {
	data         []byte
	cs           int
	p, pe, eof   int
	pb           int
	err          error
	bestEffort   bool
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
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

// Err returns the last error occurred.
//
// If the result is nil, then the parsing was successfull.
func (m *machine) Err() error {
	return m.err
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

	%% write init;
	%% write exec;

	if m.cs < first_final || m.cs == en_fail {
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