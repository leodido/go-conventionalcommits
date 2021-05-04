package parser

import (
	"bytes"
	"fmt"

	"github.com/leodido/go-conventionalcommits"
	"github.com/sirupsen/logrus"
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
	// ErrMalformedScopeClosing represents a specific early-exit error.
	ErrMalformedScopeClosing = "expecting closing parentheses (')') character, got early exit at '%s' character"
	// ErrDescriptionInit tells the user that before of the description part a whitespace is mandatory.
	ErrDescriptionInit = "expecting at least one white-space (' ') character, got '%s' character"
	// ErrDescription tells the user that after the whitespace is mandatory a description.
	ErrDescription = "expecting a description text (without newlines) after '%s' character"
	// ErrNewline communicates an illegal newline to the user.
	ErrNewline = "illegal newline"
	// ErrMissingBlankLineAtBeginning tells the user that the a blank line is missing after the description or after the body.
	ErrMissingBlankLineAtBeginning = "missing a blank line"
)

const start int = 1
const firstFinal int = 110

const enTrailerBeg int = 112
const enTrailerEnd int = 18
const enBody int = 19
const enMain int = 1
const enConventionalTypesMain int = 20
const enFalcoTypesMain int = 61
const enFreeFormTypesMain int = 101

type machine struct {
	data             []byte
	cs               int
	p, pe, eof       int
	pb               int
	err              error
	bestEffort       bool
	typeConfig       conventionalcommits.TypeConfig
	logger           *logrus.Logger
	currentFooterKey string
	countNewlines    int
	lastNewline      int
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

func (m *machine) emitInfo(s string, args ...interface{}) {
	if m.logger != nil {
		logEntry := logrus.NewEntry(m.logger)
		for i := 0; i < len(args); i = i + 2 {
			logEntry = m.logger.WithField(args[0].(string), args[1])
		}
		logEntry.Infoln(s)
	}
}

func (m *machine) emitDebug(s string, args ...interface{}) {
	if m.logger != nil {
		logEntry := logrus.NewEntry(m.logger)
		for i := 0; i < len(args); i = i + 2 {
			logEntry = m.logger.WithField(args[0].(string), args[1])
		}
		logEntry.Debugln(s)
	}
}

func (m *machine) emitError(s string, args ...interface{}) error {
	e := fmt.Errorf(s+ColumnPositionTemplate, args...)
	if m.logger != nil {
		m.logger.Errorln(e)
	}
	return e
}

func (m *machine) emitErrorWithoutCharacter(messageTemplate string) error {
	return m.emitError(messageTemplate, m.p)
}

func (m *machine) emitErrorOnCurrentCharacter(messageTemplate string) error {
	return m.emitError(messageTemplate, string(m.data[m.p]), m.p)
}

func (m *machine) emitErrorOnPreviousCharacter(messageTemplate string) error {
	return m.emitError(messageTemplate, string(m.data[m.p-1]), m.p)
}

// NewMachine creates a new FSM able to parse Conventional Commits.
func NewMachine(options ...conventionalcommits.MachineOption) conventionalcommits.Machine {
	m := &machine{}

	for _, opt := range options {
		opt(m)
	}

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
	m.currentFooterKey = ""
	m.countNewlines = 0
	output := &conventionalCommit{}
	output.footers = make(map[string][]string)

	switch m.typeConfig {
	case conventionalcommits.TypesFreeForm:
		m.cs = enFreeFormTypesMain
		break
	case conventionalcommits.TypesFalco:
		m.cs = enFalcoTypesMain
		break
	case conventionalcommits.TypesConventional:
		m.cs = enConventionalTypesMain
		break
	case conventionalcommits.TypesMinimal:
		fallthrough
	default:

		{
			m.cs = start
		}

		break
	}

	{
		var _widec int16
		if (m.p) == (m.pe) {
			goto _testEof
		}
		switch m.cs {
		case 1:
			goto stCase1
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
		case 6:
			goto stCase6
		case 7:
			goto stCase7
		case 8:
			goto stCase8
		case 110:
			goto stCase110
		case 9:
			goto stCase9
		case 111:
			goto stCase111
		case 10:
			goto stCase10
		case 11:
			goto stCase11
		case 12:
			goto stCase12
		case 13:
			goto stCase13
		case 18:
			goto stCase18
		case 114:
			goto stCase114
		case 115:
			goto stCase115
		case 19:
			goto stCase19
		case 116:
			goto stCase116
		case 20:
			goto stCase20
		case 21:
			goto stCase21
		case 22:
			goto stCase22
		case 23:
			goto stCase23
		case 24:
			goto stCase24
		case 25:
			goto stCase25
		case 26:
			goto stCase26
		case 27:
			goto stCase27
		case 28:
			goto stCase28
		case 117:
			goto stCase117
		case 29:
			goto stCase29
		case 118:
			goto stCase118
		case 30:
			goto stCase30
		case 31:
			goto stCase31
		case 32:
			goto stCase32
		case 33:
			goto stCase33
		case 34:
			goto stCase34
		case 35:
			goto stCase35
		case 36:
			goto stCase36
		case 37:
			goto stCase37
		case 38:
			goto stCase38
		case 39:
			goto stCase39
		case 40:
			goto stCase40
		case 41:
			goto stCase41
		case 42:
			goto stCase42
		case 43:
			goto stCase43
		case 44:
			goto stCase44
		case 45:
			goto stCase45
		case 46:
			goto stCase46
		case 47:
			goto stCase47
		case 48:
			goto stCase48
		case 49:
			goto stCase49
		case 50:
			goto stCase50
		case 51:
			goto stCase51
		case 52:
			goto stCase52
		case 53:
			goto stCase53
		case 54:
			goto stCase54
		case 55:
			goto stCase55
		case 56:
			goto stCase56
		case 57:
			goto stCase57
		case 58:
			goto stCase58
		case 59:
			goto stCase59
		case 60:
			goto stCase60
		case 61:
			goto stCase61
		case 62:
			goto stCase62
		case 63:
			goto stCase63
		case 64:
			goto stCase64
		case 65:
			goto stCase65
		case 66:
			goto stCase66
		case 67:
			goto stCase67
		case 68:
			goto stCase68
		case 69:
			goto stCase69
		case 119:
			goto stCase119
		case 70:
			goto stCase70
		case 120:
			goto stCase120
		case 71:
			goto stCase71
		case 72:
			goto stCase72
		case 73:
			goto stCase73
		case 74:
			goto stCase74
		case 75:
			goto stCase75
		case 76:
			goto stCase76
		case 77:
			goto stCase77
		case 78:
			goto stCase78
		case 79:
			goto stCase79
		case 80:
			goto stCase80
		case 81:
			goto stCase81
		case 82:
			goto stCase82
		case 83:
			goto stCase83
		case 84:
			goto stCase84
		case 85:
			goto stCase85
		case 86:
			goto stCase86
		case 87:
			goto stCase87
		case 88:
			goto stCase88
		case 89:
			goto stCase89
		case 90:
			goto stCase90
		case 91:
			goto stCase91
		case 92:
			goto stCase92
		case 93:
			goto stCase93
		case 94:
			goto stCase94
		case 95:
			goto stCase95
		case 96:
			goto stCase96
		case 97:
			goto stCase97
		case 98:
			goto stCase98
		case 99:
			goto stCase99
		case 100:
			goto stCase100
		case 101:
			goto stCase101
		case 102:
			goto stCase102
		case 103:
			goto stCase103
		case 104:
			goto stCase104
		case 105:
			goto stCase105
		case 121:
			goto stCase121
		case 106:
			goto stCase106
		case 122:
			goto stCase122
		case 107:
			goto stCase107
		case 108:
			goto stCase108
		case 109:
			goto stCase109
		case 112:
			goto stCase112
		case 14:
			goto stCase14
		case 15:
			goto stCase15
		case 113:
			goto stCase113
		case 16:
			goto stCase16
		case 17:
			goto stCase17
		}
		goto stOut
	stCase1:
		switch (m.data)[(m.p)] {
		case 70:
			goto tr1
		case 102:
			goto tr1
		}
		goto tr0
	tr0:

		if m.pe > 0 {
			if m.p != m.pe {
				m.err = m.emitErrorOnCurrentCharacter(ErrType)
			} else {
				m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
			}
		}

		goto st0
	tr6:

		if m.err == nil {
			m.err = m.emitErrorOnCurrentCharacter(ErrColon)
		}

		goto st0
	tr10:

		if m.err == nil {
			m.err = m.emitErrorOnCurrentCharacter(ErrDescriptionInit)
		}

		goto st0
	tr13:

		if m.p < m.pe && m.data[m.p] == 10 {
			m.err = m.emitError(ErrNewline, m.p+1)
		} else {
			m.err = m.emitErrorOnPreviousCharacter(ErrDescription)
		}

		goto st0
	tr14:

		m.err = m.emitErrorWithoutCharacter(ErrMissingBlankLineAtBeginning)

		goto st0
	tr16:

		if m.p < m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrMalformedScope)
		}

		goto st0
	tr21:

		if len(output.footers) == 0 {
			// Backtrack to the last marker
			// Ie., the text possibly a trailer token that is instead part of the body content
			if m.countNewlines > 0 {
				// In case new lines met while rewinding
				// advance the last marker by the number of the newlines so that they don't get parsed again
				// (they be added in the result by the body content appender)
				m.pb = m.lastNewline + 1
			}
			(m.p) = (m.pb) - 1

			m.emitDebug("try to parse body content", "pos", m.p)
			{
				goto st19
			}
		} else {
			fmt.Println("todo > rewind/continue to parse footer trailers", m.pb, m.p, string(m.text()))
		}

		goto st0
	tr29:

		// Append newlines
		for m.countNewlines > 0 {
			output.body += "\n"
			m.countNewlines--
			m.emitInfo("valid commit message body content", "body", "\n")
		}
		// Append content to body
		if m.p > m.pb {
			output.body += string(m.text())
			m.emitInfo("valid commit message body content", "body", string(m.text()))
		} else {
			// assert(m.p == m.pb)
			output.body += string(m.data[m.pb : m.pb+1])
			m.emitInfo("valid commit message body content", "body", string(m.data[m.pb:m.pb+1]))
		}

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st0
	tr135:

		// Append newlines
		for m.countNewlines > 0 {
			output.body += "\n"
			m.countNewlines--
			m.emitInfo("valid commit message body content", "body", "\n")
		}
		// Append content to body
		if m.p > m.pb {
			output.body += string(m.text())
			m.emitInfo("valid commit message body content", "body", string(m.text()))
		} else {
			// assert(m.p == m.pb)
			output.body += string(m.data[m.pb : m.pb+1])
			m.emitInfo("valid commit message body content", "body", string(m.data[m.pb:m.pb+1]))
		}

		// Append newlines
		for m.countNewlines > 0 {
			output.body += "\n"
			m.countNewlines--
			m.emitInfo("valid commit message body content", "body", "\n")
		}
		// Append content to body
		m.pb++
		m.p++
		output.body += string(m.text())
		m.emitInfo("valid commit message body content", "body", string(m.text()))
		// Do not advance over the current char
		(m.p)--

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st0
	stCase0:
	st0:
		m.cs = 0
		goto _out
	tr1:

		m.pb = m.p

		goto st2
	st2:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof2
		}
	stCase2:
		switch (m.data)[(m.p)] {
		case 69:
			goto st3
		case 73:
			goto st13
		case 101:
			goto st3
		case 105:
			goto st13
		}
		goto tr0
	st3:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof3
		}
	stCase3:
		switch (m.data)[(m.p)] {
		case 65:
			goto st4
		case 97:
			goto st4
		}
		goto tr0
	st4:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof4
		}
	stCase4:
		switch (m.data)[(m.p)] {
		case 84:
			goto st5
		case 116:
			goto st5
		}
		goto tr0
	st5:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof5
		}
	stCase5:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr7
		case 40:
			goto st10
		case 58:
			goto st7
		}
		goto tr6
	tr7:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st6
	st6:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof6
		}
	stCase6:
		if (m.data)[(m.p)] == 58 {
			goto st7
		}
		goto tr6
	st7:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof7
		}
	stCase7:
		if (m.data)[(m.p)] == 32 {
			goto st8
		}
		goto tr10
	st8:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof8
		}
	stCase8:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st8
		}
		goto tr12
	tr12:

		m.pb = m.p

		goto st110
	st110:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof110
		}
	stCase110:
		if (m.data)[(m.p)] == 10 {
			goto tr129
		}
		goto st110
	tr129:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st9
	st9:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof9
		}
	stCase9:
		if (m.data)[(m.p)] == 10 {
			goto tr15
		}
		goto tr14
	tr15:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st111
	st111:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof111
		}
	stCase111:
		goto st0
	st10:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof10
		}
	stCase10:
		if (m.data)[(m.p)] == 41 {
			goto tr18
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr17
			}
		case (m.data)[(m.p)] >= 32:
			goto tr17
		}
		goto tr16
	tr17:

		m.pb = m.p

		goto st11
	st11:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof11
		}
	stCase11:
		if (m.data)[(m.p)] == 41 {
			goto tr20
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st11
			}
		case (m.data)[(m.p)] >= 32:
			goto st11
		}
		goto tr16
	tr18:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st12
	tr20:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st12
	st12:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof12
		}
	stCase12:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr7
		case 58:
			goto st7
		}
		goto tr6
	st13:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof13
		}
	stCase13:
		switch (m.data)[(m.p)] {
		case 88:
			goto st5
		case 120:
			goto st5
		}
		goto tr0
	st18:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof18
		}
	stCase18:
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto tr27
		}
		goto st0
	tr27:

		m.pb = m.p

		goto st114
	st114:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof114
		}
	stCase114:
		if (m.data)[(m.p)] == 10 {
			goto tr132
		}
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto st114
		}
		goto st0
	tr132:

		output.footers[m.currentFooterKey] = append(output.footers[m.currentFooterKey], string(m.text()))
		m.emitInfo("valid commit message footer trailer", m.currentFooterKey, string(m.text()))

		// Increment number of newlines to use in case we're still in the body
		m.countNewlines++
		m.lastNewline = m.p
		m.emitDebug("found a newline", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st115
	tr134:

		// Increment number of newlines to use in case we're still in the body
		m.countNewlines++
		m.lastNewline = m.p
		m.emitDebug("found a newline", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st115
	st115:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof115
		}
	stCase115:
		if (m.data)[(m.p)] == 10 {
			goto tr134
		}
		goto st0
	st19:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof19
		}
	stCase19:
		_widec = int16((m.data)[(m.p)])
		_widec = 256 + (int16((m.data)[(m.p)]) - 0)
		if m.p+2 < m.pe && m.data[m.p+1] == 10 && m.data[m.p+2] == 10 {
			_widec += 256
		}
		if 256 <= _widec && _widec <= 511 {
			goto tr30
		}
		goto tr29
	tr30:

		m.pb = m.p

		goto st116
	tr136:

		// Append newlines
		for m.countNewlines > 0 {
			output.body += "\n"
			m.countNewlines--
			m.emitInfo("valid commit message body content", "body", "\n")
		}
		// Append body content
		output.body += string(m.text())
		m.emitInfo("valid commit message body content", "body", string(m.text()))

		m.pb = m.p

		goto st116
	st116:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof116
		}
	stCase116:
		_widec = int16((m.data)[(m.p)])
		_widec = 256 + (int16((m.data)[(m.p)]) - 0)
		if m.p+2 < m.pe && m.data[m.p+1] == 10 && m.data[m.p+2] == 10 {
			_widec += 256
		}
		if 256 <= _widec && _widec <= 511 {
			goto tr136
		}
		goto tr135
	stCase20:
		switch (m.data)[(m.p)] {
		case 66:
			goto tr31
		case 67:
			goto tr32
		case 68:
			goto tr33
		case 70:
			goto tr34
		case 80:
			goto tr35
		case 82:
			goto tr36
		case 83:
			goto tr37
		case 84:
			goto tr38
		case 98:
			goto tr31
		case 99:
			goto tr32
		case 100:
			goto tr33
		case 102:
			goto tr34
		case 112:
			goto tr35
		case 114:
			goto tr36
		case 115:
			goto tr37
		case 116:
			goto tr38
		}
		goto tr0
	tr31:

		m.pb = m.p

		goto st21
	st21:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof21
		}
	stCase21:
		switch (m.data)[(m.p)] {
		case 85:
			goto st22
		case 117:
			goto st22
		}
		goto tr0
	st22:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof22
		}
	stCase22:
		switch (m.data)[(m.p)] {
		case 73:
			goto st23
		case 105:
			goto st23
		}
		goto tr0
	st23:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof23
		}
	stCase23:
		switch (m.data)[(m.p)] {
		case 76:
			goto st24
		case 108:
			goto st24
		}
		goto tr0
	st24:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof24
		}
	stCase24:
		switch (m.data)[(m.p)] {
		case 68:
			goto st25
		case 100:
			goto st25
		}
		goto tr0
	st25:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof25
		}
	stCase25:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr43
		case 40:
			goto st30
		case 58:
			goto st27
		}
		goto tr6
	tr43:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st26
	st26:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof26
		}
	stCase26:
		if (m.data)[(m.p)] == 58 {
			goto st27
		}
		goto tr6
	st27:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof27
		}
	stCase27:
		if (m.data)[(m.p)] == 32 {
			goto st28
		}
		goto tr10
	st28:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof28
		}
	stCase28:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st28
		}
		goto tr47
	tr47:

		m.pb = m.p

		goto st117
	st117:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof117
		}
	stCase117:
		if (m.data)[(m.p)] == 10 {
			goto tr138
		}
		goto st117
	tr138:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st29
	st29:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof29
		}
	stCase29:
		if (m.data)[(m.p)] == 10 {
			goto tr48
		}
		goto tr14
	tr48:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st118
	st118:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof118
		}
	stCase118:
		goto st0
	st30:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof30
		}
	stCase30:
		if (m.data)[(m.p)] == 41 {
			goto tr50
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr49
			}
		case (m.data)[(m.p)] >= 32:
			goto tr49
		}
		goto tr16
	tr49:

		m.pb = m.p

		goto st31
	st31:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof31
		}
	stCase31:
		if (m.data)[(m.p)] == 41 {
			goto tr52
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st31
			}
		case (m.data)[(m.p)] >= 32:
			goto st31
		}
		goto tr16
	tr50:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st32
	tr52:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st32
	st32:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof32
		}
	stCase32:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr43
		case 58:
			goto st27
		}
		goto tr6
	tr32:

		m.pb = m.p

		goto st33
	st33:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof33
		}
	stCase33:
		switch (m.data)[(m.p)] {
		case 72:
			goto st34
		case 73:
			goto st25
		case 104:
			goto st34
		case 105:
			goto st25
		}
		goto tr0
	st34:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof34
		}
	stCase34:
		switch (m.data)[(m.p)] {
		case 79:
			goto st35
		case 111:
			goto st35
		}
		goto tr0
	st35:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof35
		}
	stCase35:
		switch (m.data)[(m.p)] {
		case 82:
			goto st36
		case 114:
			goto st36
		}
		goto tr0
	st36:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof36
		}
	stCase36:
		switch (m.data)[(m.p)] {
		case 69:
			goto st25
		case 101:
			goto st25
		}
		goto tr0
	tr33:

		m.pb = m.p

		goto st37
	st37:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof37
		}
	stCase37:
		switch (m.data)[(m.p)] {
		case 79:
			goto st38
		case 111:
			goto st38
		}
		goto tr0
	st38:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof38
		}
	stCase38:
		switch (m.data)[(m.p)] {
		case 67:
			goto st39
		case 99:
			goto st39
		}
		goto tr0
	st39:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof39
		}
	stCase39:
		switch (m.data)[(m.p)] {
		case 83:
			goto st25
		case 115:
			goto st25
		}
		goto tr0
	tr34:

		m.pb = m.p

		goto st40
	st40:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof40
		}
	stCase40:
		switch (m.data)[(m.p)] {
		case 69:
			goto st41
		case 73:
			goto st43
		case 101:
			goto st41
		case 105:
			goto st43
		}
		goto tr0
	st41:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof41
		}
	stCase41:
		switch (m.data)[(m.p)] {
		case 65:
			goto st42
		case 97:
			goto st42
		}
		goto tr0
	st42:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof42
		}
	stCase42:
		switch (m.data)[(m.p)] {
		case 84:
			goto st25
		case 116:
			goto st25
		}
		goto tr0
	st43:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof43
		}
	stCase43:
		switch (m.data)[(m.p)] {
		case 88:
			goto st25
		case 120:
			goto st25
		}
		goto tr0
	tr35:

		m.pb = m.p

		goto st44
	st44:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof44
		}
	stCase44:
		switch (m.data)[(m.p)] {
		case 69:
			goto st45
		case 101:
			goto st45
		}
		goto tr0
	st45:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof45
		}
	stCase45:
		switch (m.data)[(m.p)] {
		case 82:
			goto st46
		case 114:
			goto st46
		}
		goto tr0
	st46:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof46
		}
	stCase46:
		switch (m.data)[(m.p)] {
		case 70:
			goto st25
		case 102:
			goto st25
		}
		goto tr0
	tr36:

		m.pb = m.p

		goto st47
	st47:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof47
		}
	stCase47:
		switch (m.data)[(m.p)] {
		case 69:
			goto st48
		case 101:
			goto st48
		}
		goto tr0
	st48:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof48
		}
	stCase48:
		switch (m.data)[(m.p)] {
		case 70:
			goto st49
		case 86:
			goto st54
		case 102:
			goto st49
		case 118:
			goto st54
		}
		goto tr0
	st49:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof49
		}
	stCase49:
		switch (m.data)[(m.p)] {
		case 65:
			goto st50
		case 97:
			goto st50
		}
		goto tr0
	st50:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof50
		}
	stCase50:
		switch (m.data)[(m.p)] {
		case 67:
			goto st51
		case 99:
			goto st51
		}
		goto tr0
	st51:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof51
		}
	stCase51:
		switch (m.data)[(m.p)] {
		case 84:
			goto st52
		case 116:
			goto st52
		}
		goto tr0
	st52:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof52
		}
	stCase52:
		switch (m.data)[(m.p)] {
		case 79:
			goto st53
		case 111:
			goto st53
		}
		goto tr0
	st53:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof53
		}
	stCase53:
		switch (m.data)[(m.p)] {
		case 82:
			goto st25
		case 114:
			goto st25
		}
		goto tr0
	st54:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof54
		}
	stCase54:
		switch (m.data)[(m.p)] {
		case 69:
			goto st55
		case 101:
			goto st55
		}
		goto tr0
	st55:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof55
		}
	stCase55:
		switch (m.data)[(m.p)] {
		case 82:
			goto st42
		case 114:
			goto st42
		}
		goto tr0
	tr37:

		m.pb = m.p

		goto st56
	st56:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof56
		}
	stCase56:
		switch (m.data)[(m.p)] {
		case 84:
			goto st57
		case 116:
			goto st57
		}
		goto tr0
	st57:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof57
		}
	stCase57:
		switch (m.data)[(m.p)] {
		case 89:
			goto st58
		case 121:
			goto st58
		}
		goto tr0
	st58:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof58
		}
	stCase58:
		switch (m.data)[(m.p)] {
		case 76:
			goto st36
		case 108:
			goto st36
		}
		goto tr0
	tr38:

		m.pb = m.p

		goto st59
	st59:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof59
		}
	stCase59:
		switch (m.data)[(m.p)] {
		case 69:
			goto st60
		case 101:
			goto st60
		}
		goto tr0
	st60:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof60
		}
	stCase60:
		switch (m.data)[(m.p)] {
		case 83:
			goto st42
		case 115:
			goto st42
		}
		goto tr0
	stCase61:
		switch (m.data)[(m.p)] {
		case 66:
			goto tr74
		case 67:
			goto tr75
		case 68:
			goto tr76
		case 70:
			goto tr77
		case 78:
			goto tr78
		case 80:
			goto tr79
		case 82:
			goto tr80
		case 84:
			goto tr81
		case 85:
			goto tr82
		case 98:
			goto tr74
		case 99:
			goto tr75
		case 100:
			goto tr76
		case 102:
			goto tr77
		case 110:
			goto tr78
		case 112:
			goto tr79
		case 114:
			goto tr80
		case 116:
			goto tr81
		case 117:
			goto tr82
		}
		goto tr0
	tr74:

		m.pb = m.p

		goto st62
	st62:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof62
		}
	stCase62:
		switch (m.data)[(m.p)] {
		case 85:
			goto st63
		case 117:
			goto st63
		}
		goto tr0
	st63:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof63
		}
	stCase63:
		switch (m.data)[(m.p)] {
		case 73:
			goto st64
		case 105:
			goto st64
		}
		goto tr0
	st64:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof64
		}
	stCase64:
		switch (m.data)[(m.p)] {
		case 76:
			goto st65
		case 108:
			goto st65
		}
		goto tr0
	st65:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof65
		}
	stCase65:
		switch (m.data)[(m.p)] {
		case 68:
			goto st66
		case 100:
			goto st66
		}
		goto tr0
	st66:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof66
		}
	stCase66:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr87
		case 40:
			goto st71
		case 58:
			goto st68
		}
		goto tr6
	tr87:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st67
	st67:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof67
		}
	stCase67:
		if (m.data)[(m.p)] == 58 {
			goto st68
		}
		goto tr6
	st68:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof68
		}
	stCase68:
		if (m.data)[(m.p)] == 32 {
			goto st69
		}
		goto tr10
	st69:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof69
		}
	stCase69:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st69
		}
		goto tr91
	tr91:

		m.pb = m.p

		goto st119
	st119:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof119
		}
	stCase119:
		if (m.data)[(m.p)] == 10 {
			goto tr140
		}
		goto st119
	tr140:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st70
	st70:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof70
		}
	stCase70:
		if (m.data)[(m.p)] == 10 {
			goto tr92
		}
		goto tr14
	tr92:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st120
	st120:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof120
		}
	stCase120:
		goto st0
	st71:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof71
		}
	stCase71:
		if (m.data)[(m.p)] == 41 {
			goto tr94
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr93
			}
		case (m.data)[(m.p)] >= 32:
			goto tr93
		}
		goto tr16
	tr93:

		m.pb = m.p

		goto st72
	st72:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof72
		}
	stCase72:
		if (m.data)[(m.p)] == 41 {
			goto tr96
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st72
			}
		case (m.data)[(m.p)] >= 32:
			goto st72
		}
		goto tr16
	tr94:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st73
	tr96:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st73
	st73:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof73
		}
	stCase73:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr87
		case 58:
			goto st68
		}
		goto tr6
	tr75:

		m.pb = m.p

		goto st74
	st74:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof74
		}
	stCase74:
		switch (m.data)[(m.p)] {
		case 72:
			goto st75
		case 73:
			goto st66
		case 104:
			goto st75
		case 105:
			goto st66
		}
		goto tr0
	st75:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof75
		}
	stCase75:
		switch (m.data)[(m.p)] {
		case 79:
			goto st76
		case 111:
			goto st76
		}
		goto tr0
	st76:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof76
		}
	stCase76:
		switch (m.data)[(m.p)] {
		case 82:
			goto st77
		case 114:
			goto st77
		}
		goto tr0
	st77:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof77
		}
	stCase77:
		switch (m.data)[(m.p)] {
		case 69:
			goto st66
		case 101:
			goto st66
		}
		goto tr0
	tr76:

		m.pb = m.p

		goto st78
	st78:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof78
		}
	stCase78:
		switch (m.data)[(m.p)] {
		case 79:
			goto st79
		case 111:
			goto st79
		}
		goto tr0
	st79:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof79
		}
	stCase79:
		switch (m.data)[(m.p)] {
		case 67:
			goto st80
		case 99:
			goto st80
		}
		goto tr0
	st80:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof80
		}
	stCase80:
		switch (m.data)[(m.p)] {
		case 83:
			goto st66
		case 115:
			goto st66
		}
		goto tr0
	tr77:

		m.pb = m.p

		goto st81
	st81:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof81
		}
	stCase81:
		switch (m.data)[(m.p)] {
		case 69:
			goto st82
		case 73:
			goto st84
		case 101:
			goto st82
		case 105:
			goto st84
		}
		goto tr0
	st82:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof82
		}
	stCase82:
		switch (m.data)[(m.p)] {
		case 65:
			goto st83
		case 97:
			goto st83
		}
		goto tr0
	st83:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof83
		}
	stCase83:
		switch (m.data)[(m.p)] {
		case 84:
			goto st66
		case 116:
			goto st66
		}
		goto tr0
	st84:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof84
		}
	stCase84:
		switch (m.data)[(m.p)] {
		case 88:
			goto st66
		case 120:
			goto st66
		}
		goto tr0
	tr78:

		m.pb = m.p

		goto st85
	st85:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof85
		}
	stCase85:
		switch (m.data)[(m.p)] {
		case 69:
			goto st86
		case 101:
			goto st86
		}
		goto tr0
	st86:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof86
		}
	stCase86:
		switch (m.data)[(m.p)] {
		case 87:
			goto st66
		case 119:
			goto st66
		}
		goto tr0
	tr79:

		m.pb = m.p

		goto st87
	st87:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof87
		}
	stCase87:
		switch (m.data)[(m.p)] {
		case 69:
			goto st88
		case 101:
			goto st88
		}
		goto tr0
	st88:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof88
		}
	stCase88:
		switch (m.data)[(m.p)] {
		case 82:
			goto st89
		case 114:
			goto st89
		}
		goto tr0
	st89:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof89
		}
	stCase89:
		switch (m.data)[(m.p)] {
		case 70:
			goto st66
		case 102:
			goto st66
		}
		goto tr0
	tr80:

		m.pb = m.p

		goto st90
	st90:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof90
		}
	stCase90:
		switch (m.data)[(m.p)] {
		case 69:
			goto st91
		case 85:
			goto st94
		case 101:
			goto st91
		case 117:
			goto st94
		}
		goto tr0
	st91:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof91
		}
	stCase91:
		switch (m.data)[(m.p)] {
		case 86:
			goto st92
		case 118:
			goto st92
		}
		goto tr0
	st92:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof92
		}
	stCase92:
		switch (m.data)[(m.p)] {
		case 69:
			goto st93
		case 101:
			goto st93
		}
		goto tr0
	st93:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof93
		}
	stCase93:
		switch (m.data)[(m.p)] {
		case 82:
			goto st83
		case 114:
			goto st83
		}
		goto tr0
	st94:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof94
		}
	stCase94:
		switch (m.data)[(m.p)] {
		case 76:
			goto st77
		case 108:
			goto st77
		}
		goto tr0
	tr81:

		m.pb = m.p

		goto st95
	st95:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof95
		}
	stCase95:
		switch (m.data)[(m.p)] {
		case 69:
			goto st96
		case 101:
			goto st96
		}
		goto tr0
	st96:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof96
		}
	stCase96:
		switch (m.data)[(m.p)] {
		case 83:
			goto st83
		case 115:
			goto st83
		}
		goto tr0
	tr82:

		m.pb = m.p

		goto st97
	st97:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof97
		}
	stCase97:
		switch (m.data)[(m.p)] {
		case 80:
			goto st98
		case 112:
			goto st98
		}
		goto tr0
	st98:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof98
		}
	stCase98:
		switch (m.data)[(m.p)] {
		case 68:
			goto st99
		case 100:
			goto st99
		}
		goto tr0
	st99:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof99
		}
	stCase99:
		switch (m.data)[(m.p)] {
		case 65:
			goto st100
		case 97:
			goto st100
		}
		goto tr0
	st100:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof100
		}
	stCase100:
		switch (m.data)[(m.p)] {
		case 84:
			goto st77
		case 116:
			goto st77
		}
		goto tr0
	stCase101:
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto tr116
		}
		goto tr0
	tr116:

		m.pb = m.p

		goto st102
	st102:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof102
		}
	stCase102:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr118
		case 40:
			goto st107
		case 58:
			goto st104
		}
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto st102
		}
		goto tr6
	tr118:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st103
	st103:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof103
		}
	stCase103:
		if (m.data)[(m.p)] == 58 {
			goto st104
		}
		goto tr6
	st104:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof104
		}
	stCase104:
		if (m.data)[(m.p)] == 32 {
			goto st105
		}
		goto tr10
	st105:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof105
		}
	stCase105:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st105
		}
		goto tr122
	tr122:

		m.pb = m.p

		goto st121
	st121:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof121
		}
	stCase121:
		if (m.data)[(m.p)] == 10 {
			goto tr142
		}
		goto st121
	tr142:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st106
	st106:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof106
		}
	stCase106:
		if (m.data)[(m.p)] == 10 {
			goto tr123
		}
		goto tr14
	tr123:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st112
		}

		goto st122
	st122:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof122
		}
	stCase122:
		goto st0
	st107:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof107
		}
	stCase107:
		if (m.data)[(m.p)] == 41 {
			goto tr125
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr124
			}
		case (m.data)[(m.p)] >= 32:
			goto tr124
		}
		goto tr16
	tr124:

		m.pb = m.p

		goto st108
	st108:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof108
		}
	stCase108:
		if (m.data)[(m.p)] == 41 {
			goto tr127
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st108
			}
		case (m.data)[(m.p)] >= 32:
			goto st108
		}
		goto tr16
	tr125:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st109
	tr127:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st109
	st109:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof109
		}
	stCase109:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr118
		case 58:
			goto st104
		}
		goto tr6
	tr130:

		// Increment number of newlines to use in case we're still in the body
		m.countNewlines++
		m.lastNewline = m.p
		m.emitDebug("found a newline", "pos", m.p)

		goto st112
	st112:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof112
		}
	stCase112:
		if (m.data)[(m.p)] == 10 {
			goto tr130
		}
		switch {
		case (m.data)[(m.p)] < 65:
			if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
				goto tr131
			}
		case (m.data)[(m.p)] > 90:
			if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
				goto tr131
			}
		default:
			goto tr131
		}
		goto tr21
	tr131:

		m.pb = m.p

		goto st14
	st14:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof14
		}
	stCase14:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		}
		switch {
		case (m.data)[(m.p)] < 65:
			if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
				goto st14
			}
		case (m.data)[(m.p)] > 90:
			if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
				goto st14
			}
		default:
			goto st14
		}
		goto tr21
	tr22:

		m.currentFooterKey = string(bytes.ToLower(m.text()))
		m.emitDebug("possibly valid footer token", "token", m.currentFooterKey, "pos", m.p)

		goto st15
	st15:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof15
		}
	stCase15:
		if (m.data)[(m.p)] == 35 {
			goto tr26
		}
		goto tr21
	tr26:

		m.emitDebug("try to parse a footer trailer value", "pos", m.p)
		{
			goto st18
		}

		goto st113
	st113:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof113
		}
	stCase113:
		goto st0
	st16:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof16
		}
	stCase16:
		switch {
		case (m.data)[(m.p)] < 65:
			if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
				goto st14
			}
		case (m.data)[(m.p)] > 90:
			if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
				goto st14
			}
		default:
			goto st14
		}
		goto tr21
	tr25:

		m.currentFooterKey = string(bytes.ToLower(m.text()))
		m.emitDebug("possibly valid footer token", "token", m.currentFooterKey, "pos", m.p)

		goto st17
	st17:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof17
		}
	stCase17:
		if (m.data)[(m.p)] == 32 {
			goto tr26
		}
		goto tr21
	stOut:
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
	_testEof6:
		m.cs = 6
		goto _testEof
	_testEof7:
		m.cs = 7
		goto _testEof
	_testEof8:
		m.cs = 8
		goto _testEof
	_testEof110:
		m.cs = 110
		goto _testEof
	_testEof9:
		m.cs = 9
		goto _testEof
	_testEof111:
		m.cs = 111
		goto _testEof
	_testEof10:
		m.cs = 10
		goto _testEof
	_testEof11:
		m.cs = 11
		goto _testEof
	_testEof12:
		m.cs = 12
		goto _testEof
	_testEof13:
		m.cs = 13
		goto _testEof
	_testEof18:
		m.cs = 18
		goto _testEof
	_testEof114:
		m.cs = 114
		goto _testEof
	_testEof115:
		m.cs = 115
		goto _testEof
	_testEof19:
		m.cs = 19
		goto _testEof
	_testEof116:
		m.cs = 116
		goto _testEof
	_testEof21:
		m.cs = 21
		goto _testEof
	_testEof22:
		m.cs = 22
		goto _testEof
	_testEof23:
		m.cs = 23
		goto _testEof
	_testEof24:
		m.cs = 24
		goto _testEof
	_testEof25:
		m.cs = 25
		goto _testEof
	_testEof26:
		m.cs = 26
		goto _testEof
	_testEof27:
		m.cs = 27
		goto _testEof
	_testEof28:
		m.cs = 28
		goto _testEof
	_testEof117:
		m.cs = 117
		goto _testEof
	_testEof29:
		m.cs = 29
		goto _testEof
	_testEof118:
		m.cs = 118
		goto _testEof
	_testEof30:
		m.cs = 30
		goto _testEof
	_testEof31:
		m.cs = 31
		goto _testEof
	_testEof32:
		m.cs = 32
		goto _testEof
	_testEof33:
		m.cs = 33
		goto _testEof
	_testEof34:
		m.cs = 34
		goto _testEof
	_testEof35:
		m.cs = 35
		goto _testEof
	_testEof36:
		m.cs = 36
		goto _testEof
	_testEof37:
		m.cs = 37
		goto _testEof
	_testEof38:
		m.cs = 38
		goto _testEof
	_testEof39:
		m.cs = 39
		goto _testEof
	_testEof40:
		m.cs = 40
		goto _testEof
	_testEof41:
		m.cs = 41
		goto _testEof
	_testEof42:
		m.cs = 42
		goto _testEof
	_testEof43:
		m.cs = 43
		goto _testEof
	_testEof44:
		m.cs = 44
		goto _testEof
	_testEof45:
		m.cs = 45
		goto _testEof
	_testEof46:
		m.cs = 46
		goto _testEof
	_testEof47:
		m.cs = 47
		goto _testEof
	_testEof48:
		m.cs = 48
		goto _testEof
	_testEof49:
		m.cs = 49
		goto _testEof
	_testEof50:
		m.cs = 50
		goto _testEof
	_testEof51:
		m.cs = 51
		goto _testEof
	_testEof52:
		m.cs = 52
		goto _testEof
	_testEof53:
		m.cs = 53
		goto _testEof
	_testEof54:
		m.cs = 54
		goto _testEof
	_testEof55:
		m.cs = 55
		goto _testEof
	_testEof56:
		m.cs = 56
		goto _testEof
	_testEof57:
		m.cs = 57
		goto _testEof
	_testEof58:
		m.cs = 58
		goto _testEof
	_testEof59:
		m.cs = 59
		goto _testEof
	_testEof60:
		m.cs = 60
		goto _testEof
	_testEof62:
		m.cs = 62
		goto _testEof
	_testEof63:
		m.cs = 63
		goto _testEof
	_testEof64:
		m.cs = 64
		goto _testEof
	_testEof65:
		m.cs = 65
		goto _testEof
	_testEof66:
		m.cs = 66
		goto _testEof
	_testEof67:
		m.cs = 67
		goto _testEof
	_testEof68:
		m.cs = 68
		goto _testEof
	_testEof69:
		m.cs = 69
		goto _testEof
	_testEof119:
		m.cs = 119
		goto _testEof
	_testEof70:
		m.cs = 70
		goto _testEof
	_testEof120:
		m.cs = 120
		goto _testEof
	_testEof71:
		m.cs = 71
		goto _testEof
	_testEof72:
		m.cs = 72
		goto _testEof
	_testEof73:
		m.cs = 73
		goto _testEof
	_testEof74:
		m.cs = 74
		goto _testEof
	_testEof75:
		m.cs = 75
		goto _testEof
	_testEof76:
		m.cs = 76
		goto _testEof
	_testEof77:
		m.cs = 77
		goto _testEof
	_testEof78:
		m.cs = 78
		goto _testEof
	_testEof79:
		m.cs = 79
		goto _testEof
	_testEof80:
		m.cs = 80
		goto _testEof
	_testEof81:
		m.cs = 81
		goto _testEof
	_testEof82:
		m.cs = 82
		goto _testEof
	_testEof83:
		m.cs = 83
		goto _testEof
	_testEof84:
		m.cs = 84
		goto _testEof
	_testEof85:
		m.cs = 85
		goto _testEof
	_testEof86:
		m.cs = 86
		goto _testEof
	_testEof87:
		m.cs = 87
		goto _testEof
	_testEof88:
		m.cs = 88
		goto _testEof
	_testEof89:
		m.cs = 89
		goto _testEof
	_testEof90:
		m.cs = 90
		goto _testEof
	_testEof91:
		m.cs = 91
		goto _testEof
	_testEof92:
		m.cs = 92
		goto _testEof
	_testEof93:
		m.cs = 93
		goto _testEof
	_testEof94:
		m.cs = 94
		goto _testEof
	_testEof95:
		m.cs = 95
		goto _testEof
	_testEof96:
		m.cs = 96
		goto _testEof
	_testEof97:
		m.cs = 97
		goto _testEof
	_testEof98:
		m.cs = 98
		goto _testEof
	_testEof99:
		m.cs = 99
		goto _testEof
	_testEof100:
		m.cs = 100
		goto _testEof
	_testEof102:
		m.cs = 102
		goto _testEof
	_testEof103:
		m.cs = 103
		goto _testEof
	_testEof104:
		m.cs = 104
		goto _testEof
	_testEof105:
		m.cs = 105
		goto _testEof
	_testEof121:
		m.cs = 121
		goto _testEof
	_testEof106:
		m.cs = 106
		goto _testEof
	_testEof122:
		m.cs = 122
		goto _testEof
	_testEof107:
		m.cs = 107
		goto _testEof
	_testEof108:
		m.cs = 108
		goto _testEof
	_testEof109:
		m.cs = 109
		goto _testEof
	_testEof112:
		m.cs = 112
		goto _testEof
	_testEof14:
		m.cs = 14
		goto _testEof
	_testEof15:
		m.cs = 15
		goto _testEof
	_testEof113:
		m.cs = 113
		goto _testEof
	_testEof16:
		m.cs = 16
		goto _testEof
	_testEof17:
		m.cs = 17
		goto _testEof

	_testEof:
		{
		}
		if (m.p) == (m.eof) {
			switch m.cs {
			case 2, 3, 4, 13, 21, 22, 23, 24, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 62, 63, 64, 65, 74, 75, 76, 77, 78, 79, 80, 81, 82, 83, 84, 85, 86, 87, 88, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100:

				if m.pe > 0 {
					if m.p != m.pe {
						m.err = m.emitErrorOnCurrentCharacter(ErrType)
					} else {
						m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
					}
				}

			case 5, 6, 12, 25, 26, 32, 66, 67, 73, 102, 103, 109:

				if m.err == nil {
					m.err = m.emitErrorOnCurrentCharacter(ErrColon)
				}

			case 7, 27, 68, 104:

				if m.err == nil {
					m.err = m.emitErrorOnCurrentCharacter(ErrDescriptionInit)
				}

			case 8, 28, 69, 105:

				if m.p < m.pe && m.data[m.p] == 10 {
					m.err = m.emitError(ErrNewline, m.p+1)
				} else {
					m.err = m.emitErrorOnPreviousCharacter(ErrDescription)
				}

			case 9, 29, 70, 106:

				m.err = m.emitErrorWithoutCharacter(ErrMissingBlankLineAtBeginning)

			case 110, 117, 119, 121:

				output.descr = string(m.text())
				m.emitInfo("valid commit message description", "description", output.descr)

			case 114:

				output.footers[m.currentFooterKey] = append(output.footers[m.currentFooterKey], string(m.text()))
				m.emitInfo("valid commit message footer trailer", m.currentFooterKey, string(m.text()))

			case 116:

				// Append newlines
				for m.countNewlines > 0 {
					output.body += "\n"
					m.countNewlines--
					m.emitInfo("valid commit message body content", "body", "\n")
				}
				// Append body content
				output.body += string(m.text())
				m.emitInfo("valid commit message body content", "body", string(m.text()))

			case 14, 15, 16, 17:

				if len(output.footers) == 0 {
					// Backtrack to the last marker
					// Ie., the text possibly a trailer token that is instead part of the body content
					if m.countNewlines > 0 {
						// In case new lines met while rewinding
						// advance the last marker by the number of the newlines so that they don't get parsed again
						// (they be added in the result by the body content appender)
						m.pb = m.lastNewline + 1
					}
					(m.p) = (m.pb) - 1

					m.emitDebug("try to parse body content", "pos", m.p)
					{
						goto st19
					}
				} else {
					fmt.Println("todo > rewind/continue to parse footer trailers", m.pb, m.p, string(m.text()))
				}

			case 1, 20, 61, 101:

				m.err = m.emitErrorWithoutCharacter(ErrEmpty)

				if m.pe > 0 {
					if m.p != m.pe {
						m.err = m.emitErrorOnCurrentCharacter(ErrType)
					} else {
						m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
					}
				}

			case 10, 11, 30, 31, 71, 72, 107, 108:

				if m.p < m.pe {
					m.err = m.emitErrorOnCurrentCharacter(ErrMalformedScope)
				}

				// assert(m.p >= m.pe)
				m.err = m.emitErrorOnPreviousCharacter(ErrMalformedScopeClosing)

			case 19:

				// Append newlines
				for m.countNewlines > 0 {
					output.body += "\n"
					m.countNewlines--
					m.emitInfo("valid commit message body content", "body", "\n")
				}
				// Append content to body
				if m.p > m.pb {
					output.body += string(m.text())
					m.emitInfo("valid commit message body content", "body", string(m.text()))
				} else {
					// assert(m.p == m.pb)
					output.body += string(m.data[m.pb : m.pb+1])
					m.emitInfo("valid commit message body content", "body", string(m.data[m.pb:m.pb+1]))
				}

				m.emitDebug("try to parse a footer trailer token", "pos", m.p)
				{
					goto st112
				}

			}
		}

	_out:
		{
		}
	}

	if m.cs < firstFinal {
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

// WithTypes tells the parser which commit message types to consider.
func (m *machine) WithTypes(t conventionalcommits.TypeConfig) {
	m.typeConfig = t
}

// WithLogger tells the parser which logger to use.
func (m *machine) WithLogger(l *logrus.Logger) {
	m.logger = l
}
