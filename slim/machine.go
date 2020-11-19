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
	ErrColon = "missing colon (':') after '%s' character"
	// ErrTypeIncomplete represents an error in the type part of the commit message.
	ErrTypeIncomplete = "incomplete commit message type after '%s' character"
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
)

const start int = 1
const firstFinal int = 10

const enMinimalTypes int = 2
const enScope int = 12
const enBreaking int = 14
const enSeparator int = 9
const enFail int = 17
const enMain int = 1

type machine struct {
	data       []byte
	cs         int
	p, pe, eof int
	pb         int
	err        error
	bestEffort bool
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

func (m *machine) emitErrorOnCurrentCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate+ColumnPositionTemplate, string(m.data[m.p]), m.p)
}

func (m *machine) emitErrorOnPreviousCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate+ColumnPositionTemplate, string(m.data[m.p-1]), m.p)
}

// NewMachine creates a new FSM able to parse Conventional Commits.
func NewMachine(options ...conventionalcommits.MachineOption) conventionalcommits.Machine {
	m := &machine{}

	for _, opt := range options {
		opt(m)
	}

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

	{
		m.cs = start
	}

	{
		if (m.p) == (m.pe) {
			goto _testEof
		}
		goto _resume

	_again:
		switch m.cs {
		case 1:
			goto st1
		case 10:
			goto st10
		case 0:
			goto st0
		case 2:
			goto st2
		case 3:
			goto st3
		case 4:
			goto st4
		case 5:
			goto st5
		case 11:
			goto st11
		case 6:
			goto st6
		case 9:
			goto st9
		case 16:
			goto st16
		case 12:
			goto st12
		case 7:
			goto st7
		case 8:
			goto st8
		case 13:
			goto st13
		case 14:
			goto st14
		case 15:
			goto st15
		case 17:
			goto st17
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof
		}
	_resume:
		switch m.cs {
		case 1:
			goto stCase1
		case 10:
			goto stCase10
		case 0:
			goto stCase0
		case 2:
			goto stCase2
		case 3:
			goto stCase3
		case 4:
			goto stCase4
		case 5:
			goto stCase5
		case 11:
			goto stCase11
		case 6:
			goto stCase6
		case 9:
			goto stCase9
		case 16:
			goto stCase16
		case 12:
			goto stCase12
		case 7:
			goto stCase7
		case 8:
			goto stCase8
		case 13:
			goto stCase13
		case 14:
			goto stCase14
		case 15:
			goto stCase15
		case 17:
			goto stCase17
		}
		goto stOut
	st1:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof1
		}
	stCase1:
		goto tr0
	tr0:
		m.cs = 10

		(m.p)--

		m.cs = 2

		goto _again
	st10:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof10
		}
	stCase10:
		goto st0
	tr3:

		if m.p != m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrType)
		} else {
			m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
		}

		goto st0
	stCase0:
	st0:
		m.cs = 0
		goto _out
	st2:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof2
		}
	stCase2:
		if (m.data)[(m.p)] == 102 {
			goto tr1
		}
		goto st0
	tr1:

		m.pb = m.p

		goto st3
	st3:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof3
		}
	stCase3:
		switch (m.data)[(m.p)] {
		case 101:
			goto st4
		case 105:
			goto st6
		}
		goto tr3
	st4:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof4
		}
	stCase4:
		if (m.data)[(m.p)] == 97 {
			goto st5
		}
		goto tr3
	st5:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof5
		}
	stCase5:
		if (m.data)[(m.p)] == 116 {
			goto st11
		}
		goto tr3
	st11:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof11
		}
	stCase11:
		goto st0
	st6:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof6
		}
	stCase6:
		if (m.data)[(m.p)] == 120 {
			goto st11
		}
		goto tr3
	st9:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof9
		}
	stCase9:
		if (m.data)[(m.p)] == 58 {
			goto st16
		}
		goto st0
	st16:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof16
		}
	stCase16:
		goto st0
	st12:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof12
		}
	stCase12:
		if (m.data)[(m.p)] == 40 {
			goto st7
		}
		goto st0
	st7:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof7
		}
	stCase7:
		if (m.data)[(m.p)] == 41 {
			goto tr9
		}
		goto tr8
	tr8:

		m.pb = m.p

		goto st8
	st8:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof8
		}
	stCase8:
		if (m.data)[(m.p)] == 41 {
			goto tr11
		}
		goto st8
	tr9:

		m.pb = m.p

		output.scope = string(m.text())

		goto st13
	tr11:

		output.scope = string(m.text())

		goto st13
	st13:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof13
		}
	stCase13:
		if (m.data)[(m.p)] == 41 {
			goto tr11
		}
		goto st8
	st14:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof14
		}
	stCase14:
		if (m.data)[(m.p)] == 33 {
			goto tr14
		}
		goto st0
	tr14:

		m.pb = m.p

		goto st15
	st15:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof15
		}
	stCase15:
		goto st0
	st17:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof17
		}
	stCase17:
		switch (m.data)[(m.p)] {
		case 10:
			goto st0
		case 13:
			goto st0
		}
		goto st17
	stOut:
	_testEof1:
		m.cs = 1
		goto _testEof
	_testEof10:
		m.cs = 10
		goto _testEof
	_testEof2:
		m.cs = 2
		goto _testEof
	_testEof3:
		m.cs = 3
		goto _testEof
	_testEof4:
		m.cs = 4
		goto _testEof
	_testEof5:
		m.cs = 5
		goto _testEof
	_testEof11:
		m.cs = 11
		goto _testEof
	_testEof6:
		m.cs = 6
		goto _testEof
	_testEof9:
		m.cs = 9
		goto _testEof
	_testEof16:
		m.cs = 16
		goto _testEof
	_testEof12:
		m.cs = 12
		goto _testEof
	_testEof7:
		m.cs = 7
		goto _testEof
	_testEof8:
		m.cs = 8
		goto _testEof
	_testEof13:
		m.cs = 13
		goto _testEof
	_testEof14:
		m.cs = 14
		goto _testEof
	_testEof15:
		m.cs = 15
		goto _testEof
	_testEof17:
		m.cs = 17
		goto _testEof

	_testEof:
		{
		}
		if (m.p) == (m.eof) {
			switch m.cs {
			case 1:

				m.err = fmt.Errorf(ErrEmpty+ColumnPositionTemplate, m.p)
				(m.p)--

			case 3, 4, 5, 6:

				if m.p != m.pe {
					m.err = m.emitErrorOnCurrentCharacter(ErrType)
				} else {
					m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
				}

			case 9, 16:

				m.err = m.emitErrorOnPreviousCharacter(ErrColon)
				(m.p)--

				{
					goto st17
				}

			case 12, 13:

				(m.p)--

				{
					goto st14
				}

			case 14:

				(m.p)--

				{
					goto st9
				}

			case 11:

				output._type = string(m.text())

				(m.p)--

				{
					goto st12
				}

			case 15:

				output.exclamation = true

				(m.p)--

				{
					goto st9
				}

			}
		}

	_out:
		{
		}
	}

	if m.cs < firstFinal || m.cs == enFail {
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
