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
	// ErrTypeIncomplete represents an error when the type part of the commit message is not complete.
	ErrTypeIncomplete = "incomplete commit message type after '%s' character"
	// ErrColon is the error message that communicate that the mandatory colon after the type part of the commit message is missing.
	ErrColon = "expecting colon (':') character, got '%s' character"
	// ErrScope represents an error about illegal characters into the the scope part of the commit message.
	ErrScope = "illegal '%s' character in scope"
	// ErrScopeIncomplete represents a specific early-exit error.
	ErrScopeIncomplete = "expecting closing parentheses (')') character, got early exit after '%s' character"
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
	// ErrEarly represents an error when the input makes the machine exit too early.
	ErrEarly = "early exit after '%s' character"
	// ErrDescriptionInit tells the user that before of the description part a whitespace is mandatory.
	ErrDescriptionInit = "expecting at least one white-space (' ') character, got '%s' character"
	// ErrDescription tells the user that after the whitespace is mandatory a description.
	ErrDescription = "expecting a description text (without newlines) after '%s' character"
	// ErrNewline communicates an illegal newline to the user.
	ErrNewline = "illegal newline"
	// ErrMissingBlankLineAtBeginning tells the user that the a blank line is missing after the description or after the body.
	ErrMissingBlankLineAtBeginning = "missing a blank line"
	// ErrTrailer represents an error due to an unexepected character while parsing a footer trailer.
	ErrTrailer = "illegal '%s' character in trailer"
	// ErrTrailerIncomplete represent an error when a trailer is not complete.
	ErrTrailerIncomplete = "incomplete footer trailer after '%s' character"
)

const start int = 1
const firstFinal int = 125

const enTrailerBeg int = 127
const enTrailerEnd int = 33
const enBody int = 34
const enMain int = 1
const enConventionalTypesMain int = 35
const enFalcoTypesMain int = 76
const enFreeFormTypesMain int = 116

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
		case 125:
			goto stCase125
		case 9:
			goto stCase9
		case 126:
			goto stCase126
		case 10:
			goto stCase10
		case 11:
			goto stCase11
		case 12:
			goto stCase12
		case 13:
			goto stCase13
		case 33:
			goto stCase33
		case 130:
			goto stCase130
		case 131:
			goto stCase131
		case 34:
			goto stCase34
		case 132:
			goto stCase132
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
		case 133:
			goto stCase133
		case 44:
			goto stCase44
		case 134:
			goto stCase134
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
		case 70:
			goto stCase70
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
		case 135:
			goto stCase135
		case 85:
			goto stCase85
		case 136:
			goto stCase136
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
		case 106:
			goto stCase106
		case 107:
			goto stCase107
		case 108:
			goto stCase108
		case 109:
			goto stCase109
		case 110:
			goto stCase110
		case 111:
			goto stCase111
		case 112:
			goto stCase112
		case 113:
			goto stCase113
		case 114:
			goto stCase114
		case 115:
			goto stCase115
		case 116:
			goto stCase116
		case 117:
			goto stCase117
		case 118:
			goto stCase118
		case 119:
			goto stCase119
		case 120:
			goto stCase120
		case 137:
			goto stCase137
		case 121:
			goto stCase121
		case 138:
			goto stCase138
		case 122:
			goto stCase122
		case 123:
			goto stCase123
		case 124:
			goto stCase124
		case 127:
			goto stCase127
		case 14:
			goto stCase14
		case 15:
			goto stCase15
		case 128:
			goto stCase128
		case 16:
			goto stCase16
		case 17:
			goto stCase17
		case 129:
			goto stCase129
		case 18:
			goto stCase18
		case 19:
			goto stCase19
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
		case 29:
			goto stCase29
		case 30:
			goto stCase30
		case 31:
			goto stCase31
		case 32:
			goto stCase32
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
				// assert(m.p == m.pe)
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
			// assert(m.p == m.pe)
			m.err = m.emitErrorOnPreviousCharacter(ErrDescription)
		}

		goto st0
	tr14:

		m.err = m.emitErrorWithoutCharacter(ErrMissingBlankLineAtBeginning)

		goto st0
	tr16:

		if m.p < m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrScope)
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
				goto st34
			}
		} else {
			// A rewind happens when an error while parsing a footer trailer is encountered
			// If this is not the first footer trailer the parser can't go back to parse body content again
			// Thus, emit an error
			if m.p != m.pe {
				m.err = m.emitErrorOnCurrentCharacter(ErrTrailer)
			} else {
				// assert(m.p == m.pe)
				m.err = m.emitErrorOnPreviousCharacter(ErrTrailerIncomplete)
			}
		}

		goto st0
	tr44:

		// Append newlines
		for m.countNewlines > 0 {
			output.body += "\n"
			m.countNewlines--
			m.emitInfo("valid commit message body content", "body", "\n")
		}
		// Append body content
		output.body += string(m.text())
		m.emitInfo("valid commit message body content", "body", string(m.text()))

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
		}

		goto st0
	tr151:

		// Append newlines
		for m.countNewlines > 0 {
			output.body += "\n"
			m.countNewlines--
			m.emitInfo("valid commit message body content", "body", "\n")
		}
		// Append body content
		output.body += string(m.text())
		m.emitInfo("valid commit message body content", "body", string(m.text()))

		// Append content to body
		m.pb++
		m.p++
		output.body += string(m.text())
		m.emitInfo("valid commit message body content", "body", string(m.text()))
		// Do not advance over the current char
		(m.p)--

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
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

		goto st125
	st125:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof125
		}
	stCase125:
		if (m.data)[(m.p)] == 10 {
			goto tr144
		}
		goto st125
	tr144:

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
			goto st127
		}

		goto st126
	st126:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof126
		}
	stCase126:
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
	st33:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof33
		}
	stCase33:
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto tr42
		}
		goto st0
	tr42:

		m.pb = m.p

		goto st130
	st130:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof130
		}
	stCase130:
		if (m.data)[(m.p)] == 10 {
			goto tr148
		}
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto st130
		}
		goto st0
	tr148:

		output.footers[m.currentFooterKey] = append(output.footers[m.currentFooterKey], string(m.text()))
		m.emitInfo("valid commit message footer trailer", m.currentFooterKey, string(m.text()))

		// Increment number of newlines to use in case we're still in the body
		m.countNewlines++
		m.lastNewline = m.p
		m.emitDebug("found a newline", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
		}

		goto st131
	tr150:

		// Increment number of newlines to use in case we're still in the body
		m.countNewlines++
		m.lastNewline = m.p
		m.emitDebug("found a newline", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
		}

		goto st131
	st131:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof131
		}
	stCase131:
		if (m.data)[(m.p)] == 10 {
			goto tr150
		}
		goto st0
	st34:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof34
		}
	stCase34:
		_widec = int16((m.data)[(m.p)])
		_widec = 256 + (int16((m.data)[(m.p)]) - 0)
		if m.p+2 < m.pe && m.data[m.p+1] == 10 && m.data[m.p+2] == 10 {
			_widec += 256
		}
		if 256 <= _widec && _widec <= 511 {
			goto tr45
		}
		goto tr44
	tr45:

		m.pb = m.p

		goto st132
	tr152:

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

		goto st132
	st132:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof132
		}
	stCase132:
		_widec = int16((m.data)[(m.p)])
		_widec = 256 + (int16((m.data)[(m.p)]) - 0)
		if m.p+2 < m.pe && m.data[m.p+1] == 10 && m.data[m.p+2] == 10 {
			_widec += 256
		}
		if 256 <= _widec && _widec <= 511 {
			goto tr152
		}
		goto tr151
	stCase35:
		switch (m.data)[(m.p)] {
		case 66:
			goto tr46
		case 67:
			goto tr47
		case 68:
			goto tr48
		case 70:
			goto tr49
		case 80:
			goto tr50
		case 82:
			goto tr51
		case 83:
			goto tr52
		case 84:
			goto tr53
		case 98:
			goto tr46
		case 99:
			goto tr47
		case 100:
			goto tr48
		case 102:
			goto tr49
		case 112:
			goto tr50
		case 114:
			goto tr51
		case 115:
			goto tr52
		case 116:
			goto tr53
		}
		goto tr0
	tr46:

		m.pb = m.p

		goto st36
	st36:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof36
		}
	stCase36:
		switch (m.data)[(m.p)] {
		case 85:
			goto st37
		case 117:
			goto st37
		}
		goto tr0
	st37:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof37
		}
	stCase37:
		switch (m.data)[(m.p)] {
		case 73:
			goto st38
		case 105:
			goto st38
		}
		goto tr0
	st38:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof38
		}
	stCase38:
		switch (m.data)[(m.p)] {
		case 76:
			goto st39
		case 108:
			goto st39
		}
		goto tr0
	st39:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof39
		}
	stCase39:
		switch (m.data)[(m.p)] {
		case 68:
			goto st40
		case 100:
			goto st40
		}
		goto tr0
	st40:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof40
		}
	stCase40:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr58
		case 40:
			goto st45
		case 58:
			goto st42
		}
		goto tr6
	tr58:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st41
	st41:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof41
		}
	stCase41:
		if (m.data)[(m.p)] == 58 {
			goto st42
		}
		goto tr6
	st42:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof42
		}
	stCase42:
		if (m.data)[(m.p)] == 32 {
			goto st43
		}
		goto tr10
	st43:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof43
		}
	stCase43:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st43
		}
		goto tr62
	tr62:

		m.pb = m.p

		goto st133
	st133:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof133
		}
	stCase133:
		if (m.data)[(m.p)] == 10 {
			goto tr154
		}
		goto st133
	tr154:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st44
	st44:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof44
		}
	stCase44:
		if (m.data)[(m.p)] == 10 {
			goto tr63
		}
		goto tr14
	tr63:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
		}

		goto st134
	st134:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof134
		}
	stCase134:
		goto st0
	st45:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof45
		}
	stCase45:
		if (m.data)[(m.p)] == 41 {
			goto tr65
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr64
			}
		case (m.data)[(m.p)] >= 32:
			goto tr64
		}
		goto tr16
	tr64:

		m.pb = m.p

		goto st46
	st46:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof46
		}
	stCase46:
		if (m.data)[(m.p)] == 41 {
			goto tr67
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st46
			}
		case (m.data)[(m.p)] >= 32:
			goto st46
		}
		goto tr16
	tr65:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st47
	tr67:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st47
	st47:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof47
		}
	stCase47:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr58
		case 58:
			goto st42
		}
		goto tr6
	tr47:

		m.pb = m.p

		goto st48
	st48:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof48
		}
	stCase48:
		switch (m.data)[(m.p)] {
		case 72:
			goto st49
		case 73:
			goto st40
		case 104:
			goto st49
		case 105:
			goto st40
		}
		goto tr0
	st49:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof49
		}
	stCase49:
		switch (m.data)[(m.p)] {
		case 79:
			goto st50
		case 111:
			goto st50
		}
		goto tr0
	st50:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof50
		}
	stCase50:
		switch (m.data)[(m.p)] {
		case 82:
			goto st51
		case 114:
			goto st51
		}
		goto tr0
	st51:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof51
		}
	stCase51:
		switch (m.data)[(m.p)] {
		case 69:
			goto st40
		case 101:
			goto st40
		}
		goto tr0
	tr48:

		m.pb = m.p

		goto st52
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
		case 67:
			goto st54
		case 99:
			goto st54
		}
		goto tr0
	st54:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof54
		}
	stCase54:
		switch (m.data)[(m.p)] {
		case 83:
			goto st40
		case 115:
			goto st40
		}
		goto tr0
	tr49:

		m.pb = m.p

		goto st55
	st55:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof55
		}
	stCase55:
		switch (m.data)[(m.p)] {
		case 69:
			goto st56
		case 73:
			goto st58
		case 101:
			goto st56
		case 105:
			goto st58
		}
		goto tr0
	st56:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof56
		}
	stCase56:
		switch (m.data)[(m.p)] {
		case 65:
			goto st57
		case 97:
			goto st57
		}
		goto tr0
	st57:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof57
		}
	stCase57:
		switch (m.data)[(m.p)] {
		case 84:
			goto st40
		case 116:
			goto st40
		}
		goto tr0
	st58:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof58
		}
	stCase58:
		switch (m.data)[(m.p)] {
		case 88:
			goto st40
		case 120:
			goto st40
		}
		goto tr0
	tr50:

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
		case 82:
			goto st61
		case 114:
			goto st61
		}
		goto tr0
	st61:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof61
		}
	stCase61:
		switch (m.data)[(m.p)] {
		case 70:
			goto st40
		case 102:
			goto st40
		}
		goto tr0
	tr51:

		m.pb = m.p

		goto st62
	st62:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof62
		}
	stCase62:
		switch (m.data)[(m.p)] {
		case 69:
			goto st63
		case 101:
			goto st63
		}
		goto tr0
	st63:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof63
		}
	stCase63:
		switch (m.data)[(m.p)] {
		case 70:
			goto st64
		case 86:
			goto st69
		case 102:
			goto st64
		case 118:
			goto st69
		}
		goto tr0
	st64:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof64
		}
	stCase64:
		switch (m.data)[(m.p)] {
		case 65:
			goto st65
		case 97:
			goto st65
		}
		goto tr0
	st65:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof65
		}
	stCase65:
		switch (m.data)[(m.p)] {
		case 67:
			goto st66
		case 99:
			goto st66
		}
		goto tr0
	st66:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof66
		}
	stCase66:
		switch (m.data)[(m.p)] {
		case 84:
			goto st67
		case 116:
			goto st67
		}
		goto tr0
	st67:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof67
		}
	stCase67:
		switch (m.data)[(m.p)] {
		case 79:
			goto st68
		case 111:
			goto st68
		}
		goto tr0
	st68:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof68
		}
	stCase68:
		switch (m.data)[(m.p)] {
		case 82:
			goto st40
		case 114:
			goto st40
		}
		goto tr0
	st69:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof69
		}
	stCase69:
		switch (m.data)[(m.p)] {
		case 69:
			goto st70
		case 101:
			goto st70
		}
		goto tr0
	st70:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof70
		}
	stCase70:
		switch (m.data)[(m.p)] {
		case 82:
			goto st57
		case 114:
			goto st57
		}
		goto tr0
	tr52:

		m.pb = m.p

		goto st71
	st71:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof71
		}
	stCase71:
		switch (m.data)[(m.p)] {
		case 84:
			goto st72
		case 116:
			goto st72
		}
		goto tr0
	st72:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof72
		}
	stCase72:
		switch (m.data)[(m.p)] {
		case 89:
			goto st73
		case 121:
			goto st73
		}
		goto tr0
	st73:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof73
		}
	stCase73:
		switch (m.data)[(m.p)] {
		case 76:
			goto st51
		case 108:
			goto st51
		}
		goto tr0
	tr53:

		m.pb = m.p

		goto st74
	st74:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof74
		}
	stCase74:
		switch (m.data)[(m.p)] {
		case 69:
			goto st75
		case 101:
			goto st75
		}
		goto tr0
	st75:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof75
		}
	stCase75:
		switch (m.data)[(m.p)] {
		case 83:
			goto st57
		case 115:
			goto st57
		}
		goto tr0
	stCase76:
		switch (m.data)[(m.p)] {
		case 66:
			goto tr89
		case 67:
			goto tr90
		case 68:
			goto tr91
		case 70:
			goto tr92
		case 78:
			goto tr93
		case 80:
			goto tr94
		case 82:
			goto tr95
		case 84:
			goto tr96
		case 85:
			goto tr97
		case 98:
			goto tr89
		case 99:
			goto tr90
		case 100:
			goto tr91
		case 102:
			goto tr92
		case 110:
			goto tr93
		case 112:
			goto tr94
		case 114:
			goto tr95
		case 116:
			goto tr96
		case 117:
			goto tr97
		}
		goto tr0
	tr89:

		m.pb = m.p

		goto st77
	st77:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof77
		}
	stCase77:
		switch (m.data)[(m.p)] {
		case 85:
			goto st78
		case 117:
			goto st78
		}
		goto tr0
	st78:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof78
		}
	stCase78:
		switch (m.data)[(m.p)] {
		case 73:
			goto st79
		case 105:
			goto st79
		}
		goto tr0
	st79:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof79
		}
	stCase79:
		switch (m.data)[(m.p)] {
		case 76:
			goto st80
		case 108:
			goto st80
		}
		goto tr0
	st80:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof80
		}
	stCase80:
		switch (m.data)[(m.p)] {
		case 68:
			goto st81
		case 100:
			goto st81
		}
		goto tr0
	st81:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof81
		}
	stCase81:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr102
		case 40:
			goto st86
		case 58:
			goto st83
		}
		goto tr6
	tr102:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st82
	st82:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof82
		}
	stCase82:
		if (m.data)[(m.p)] == 58 {
			goto st83
		}
		goto tr6
	st83:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof83
		}
	stCase83:
		if (m.data)[(m.p)] == 32 {
			goto st84
		}
		goto tr10
	st84:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof84
		}
	stCase84:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st84
		}
		goto tr106
	tr106:

		m.pb = m.p

		goto st135
	st135:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof135
		}
	stCase135:
		if (m.data)[(m.p)] == 10 {
			goto tr156
		}
		goto st135
	tr156:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st85
	st85:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof85
		}
	stCase85:
		if (m.data)[(m.p)] == 10 {
			goto tr107
		}
		goto tr14
	tr107:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
		}

		goto st136
	st136:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof136
		}
	stCase136:
		goto st0
	st86:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof86
		}
	stCase86:
		if (m.data)[(m.p)] == 41 {
			goto tr109
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr108
			}
		case (m.data)[(m.p)] >= 32:
			goto tr108
		}
		goto tr16
	tr108:

		m.pb = m.p

		goto st87
	st87:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof87
		}
	stCase87:
		if (m.data)[(m.p)] == 41 {
			goto tr111
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st87
			}
		case (m.data)[(m.p)] >= 32:
			goto st87
		}
		goto tr16
	tr109:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st88
	tr111:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st88
	st88:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof88
		}
	stCase88:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr102
		case 58:
			goto st83
		}
		goto tr6
	tr90:

		m.pb = m.p

		goto st89
	st89:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof89
		}
	stCase89:
		switch (m.data)[(m.p)] {
		case 72:
			goto st90
		case 73:
			goto st81
		case 104:
			goto st90
		case 105:
			goto st81
		}
		goto tr0
	st90:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof90
		}
	stCase90:
		switch (m.data)[(m.p)] {
		case 79:
			goto st91
		case 111:
			goto st91
		}
		goto tr0
	st91:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof91
		}
	stCase91:
		switch (m.data)[(m.p)] {
		case 82:
			goto st92
		case 114:
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
			goto st81
		case 101:
			goto st81
		}
		goto tr0
	tr91:

		m.pb = m.p

		goto st93
	st93:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof93
		}
	stCase93:
		switch (m.data)[(m.p)] {
		case 79:
			goto st94
		case 111:
			goto st94
		}
		goto tr0
	st94:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof94
		}
	stCase94:
		switch (m.data)[(m.p)] {
		case 67:
			goto st95
		case 99:
			goto st95
		}
		goto tr0
	st95:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof95
		}
	stCase95:
		switch (m.data)[(m.p)] {
		case 83:
			goto st81
		case 115:
			goto st81
		}
		goto tr0
	tr92:

		m.pb = m.p

		goto st96
	st96:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof96
		}
	stCase96:
		switch (m.data)[(m.p)] {
		case 69:
			goto st97
		case 73:
			goto st99
		case 101:
			goto st97
		case 105:
			goto st99
		}
		goto tr0
	st97:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof97
		}
	stCase97:
		switch (m.data)[(m.p)] {
		case 65:
			goto st98
		case 97:
			goto st98
		}
		goto tr0
	st98:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof98
		}
	stCase98:
		switch (m.data)[(m.p)] {
		case 84:
			goto st81
		case 116:
			goto st81
		}
		goto tr0
	st99:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof99
		}
	stCase99:
		switch (m.data)[(m.p)] {
		case 88:
			goto st81
		case 120:
			goto st81
		}
		goto tr0
	tr93:

		m.pb = m.p

		goto st100
	st100:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof100
		}
	stCase100:
		switch (m.data)[(m.p)] {
		case 69:
			goto st101
		case 101:
			goto st101
		}
		goto tr0
	st101:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof101
		}
	stCase101:
		switch (m.data)[(m.p)] {
		case 87:
			goto st81
		case 119:
			goto st81
		}
		goto tr0
	tr94:

		m.pb = m.p

		goto st102
	st102:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof102
		}
	stCase102:
		switch (m.data)[(m.p)] {
		case 69:
			goto st103
		case 101:
			goto st103
		}
		goto tr0
	st103:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof103
		}
	stCase103:
		switch (m.data)[(m.p)] {
		case 82:
			goto st104
		case 114:
			goto st104
		}
		goto tr0
	st104:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof104
		}
	stCase104:
		switch (m.data)[(m.p)] {
		case 70:
			goto st81
		case 102:
			goto st81
		}
		goto tr0
	tr95:

		m.pb = m.p

		goto st105
	st105:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof105
		}
	stCase105:
		switch (m.data)[(m.p)] {
		case 69:
			goto st106
		case 85:
			goto st109
		case 101:
			goto st106
		case 117:
			goto st109
		}
		goto tr0
	st106:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof106
		}
	stCase106:
		switch (m.data)[(m.p)] {
		case 86:
			goto st107
		case 118:
			goto st107
		}
		goto tr0
	st107:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof107
		}
	stCase107:
		switch (m.data)[(m.p)] {
		case 69:
			goto st108
		case 101:
			goto st108
		}
		goto tr0
	st108:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof108
		}
	stCase108:
		switch (m.data)[(m.p)] {
		case 82:
			goto st98
		case 114:
			goto st98
		}
		goto tr0
	st109:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof109
		}
	stCase109:
		switch (m.data)[(m.p)] {
		case 76:
			goto st92
		case 108:
			goto st92
		}
		goto tr0
	tr96:

		m.pb = m.p

		goto st110
	st110:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof110
		}
	stCase110:
		switch (m.data)[(m.p)] {
		case 69:
			goto st111
		case 101:
			goto st111
		}
		goto tr0
	st111:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof111
		}
	stCase111:
		switch (m.data)[(m.p)] {
		case 83:
			goto st98
		case 115:
			goto st98
		}
		goto tr0
	tr97:

		m.pb = m.p

		goto st112
	st112:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof112
		}
	stCase112:
		switch (m.data)[(m.p)] {
		case 80:
			goto st113
		case 112:
			goto st113
		}
		goto tr0
	st113:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof113
		}
	stCase113:
		switch (m.data)[(m.p)] {
		case 68:
			goto st114
		case 100:
			goto st114
		}
		goto tr0
	st114:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof114
		}
	stCase114:
		switch (m.data)[(m.p)] {
		case 65:
			goto st115
		case 97:
			goto st115
		}
		goto tr0
	st115:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof115
		}
	stCase115:
		switch (m.data)[(m.p)] {
		case 84:
			goto st92
		case 116:
			goto st92
		}
		goto tr0
	stCase116:
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto tr131
		}
		goto tr0
	tr131:

		m.pb = m.p

		goto st117
	st117:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof117
		}
	stCase117:

		output._type = string(m.text())
		m.emitInfo("valid commit message type", "type", output._type)

		switch (m.data)[(m.p)] {
		case 33:
			goto tr133
		case 40:
			goto st122
		case 58:
			goto st119
		}
		if 32 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
			goto st117
		}
		goto tr6
	tr133:

		output.exclamation = true
		m.emitInfo("commit message communicates a breaking change")

		goto st118
	st118:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof118
		}
	stCase118:
		if (m.data)[(m.p)] == 58 {
			goto st119
		}
		goto tr6
	st119:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof119
		}
	stCase119:
		if (m.data)[(m.p)] == 32 {
			goto st120
		}
		goto tr10
	st120:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof120
		}
	stCase120:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr13
		case 32:
			goto st120
		}
		goto tr137
	tr137:

		m.pb = m.p

		goto st137
	st137:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof137
		}
	stCase137:
		if (m.data)[(m.p)] == 10 {
			goto tr158
		}
		goto st137
	tr158:

		output.descr = string(m.text())
		m.emitInfo("valid commit message description", "description", output.descr)

		goto st121
	st121:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof121
		}
	stCase121:
		if (m.data)[(m.p)] == 10 {
			goto tr138
		}
		goto tr14
	tr138:

		m.emitDebug("found a blank line", "pos", m.p)

		m.emitDebug("try to parse a footer trailer token", "pos", m.p)
		{
			goto st127
		}

		goto st138
	st138:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof138
		}
	stCase138:
		goto st0
	st122:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof122
		}
	stCase122:
		if (m.data)[(m.p)] == 41 {
			goto tr140
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto tr139
			}
		case (m.data)[(m.p)] >= 32:
			goto tr139
		}
		goto tr16
	tr139:

		m.pb = m.p

		goto st123
	st123:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof123
		}
	stCase123:
		if (m.data)[(m.p)] == 41 {
			goto tr142
		}
		switch {
		case (m.data)[(m.p)] > 39:
			if 42 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 126 {
				goto st123
			}
		case (m.data)[(m.p)] >= 32:
			goto st123
		}
		goto tr16
	tr140:

		m.pb = m.p

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st124
	tr142:

		output.scope = string(m.text())
		m.emitInfo("valid commit message scope", "scope", output.scope)

		goto st124
	st124:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof124
		}
	stCase124:
		switch (m.data)[(m.p)] {
		case 33:
			goto tr133
		case 58:
			goto st119
		}
		goto tr6
	tr145:

		// Increment number of newlines to use in case we're still in the body
		m.countNewlines++
		m.lastNewline = m.p
		m.emitDebug("found a newline", "pos", m.p)

		goto st127
	st127:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof127
		}
	stCase127:
		switch (m.data)[(m.p)] {
		case 10:
			goto tr145
		case 66:
			goto tr147
		}
		switch {
		case (m.data)[(m.p)] < 65:
			if 48 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 57 {
				goto tr146
			}
		case (m.data)[(m.p)] > 90:
			if 97 <= (m.data)[(m.p)] && (m.data)[(m.p)] <= 122 {
				goto tr146
			}
		default:
			goto tr146
		}
		goto tr21
	tr146:

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

		// todo > alnum[[- ]alnum] string to lower can be more performant?
		m.currentFooterKey = string(bytes.ToLower(m.text()))
		if m.currentFooterKey == "breaking change" {
			m.currentFooterKey = "breaking-change"
		}
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
			goto st33
		}

		goto st128
	st128:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof128
		}
	stCase128:
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

		// todo > alnum[[- ]alnum] string to lower can be more performant?
		m.currentFooterKey = string(bytes.ToLower(m.text()))
		if m.currentFooterKey == "breaking change" {
			m.currentFooterKey = "breaking-change"
		}
		m.emitDebug("possibly valid footer token", "token", m.currentFooterKey, "pos", m.p)

		goto st17
	st17:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof17
		}
	stCase17:
		if (m.data)[(m.p)] == 32 {
			goto tr27
		}
		goto tr21
	tr27:

		m.emitDebug("try to parse a footer trailer value", "pos", m.p)
		{
			goto st33
		}

		goto st129
	st129:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof129
		}
	stCase129:
		if (m.data)[(m.p)] == 32 {
			goto tr27
		}
		goto st0
	tr147:

		m.pb = m.p

		goto st18
	st18:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof18
		}
	stCase18:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 82:
			goto st19
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
	st19:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof19
		}
	stCase19:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 69:
			goto st20
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
	st20:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof20
		}
	stCase20:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 65:
			goto st21
		}
		switch {
		case (m.data)[(m.p)] < 66:
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
	st21:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof21
		}
	stCase21:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 75:
			goto st22
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
	st22:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof22
		}
	stCase22:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 73:
			goto st23
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
	st23:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof23
		}
	stCase23:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 78:
			goto st24
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
	st24:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof24
		}
	stCase24:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr22
		case 45:
			goto st16
		case 58:
			goto tr25
		case 71:
			goto st25
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
	st25:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof25
		}
	stCase25:
		switch (m.data)[(m.p)] {
		case 32:
			goto tr35
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
	tr35:

		// todo > alnum[[- ]alnum] string to lower can be more performant?
		m.currentFooterKey = string(bytes.ToLower(m.text()))
		if m.currentFooterKey == "breaking change" {
			m.currentFooterKey = "breaking-change"
		}
		m.emitDebug("possibly valid footer token", "token", m.currentFooterKey, "pos", m.p)

		goto st26
	st26:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof26
		}
	stCase26:
		switch (m.data)[(m.p)] {
		case 35:
			goto tr26
		case 67:
			goto st27
		}
		goto tr21
	st27:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof27
		}
	stCase27:
		if (m.data)[(m.p)] == 72 {
			goto st28
		}
		goto tr21
	st28:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof28
		}
	stCase28:
		if (m.data)[(m.p)] == 65 {
			goto st29
		}
		goto tr21
	st29:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof29
		}
	stCase29:
		if (m.data)[(m.p)] == 78 {
			goto st30
		}
		goto tr21
	st30:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof30
		}
	stCase30:
		if (m.data)[(m.p)] == 71 {
			goto st31
		}
		goto tr21
	st31:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof31
		}
	stCase31:
		if (m.data)[(m.p)] == 69 {
			goto st32
		}
		goto tr21
	st32:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof32
		}
	stCase32:
		if (m.data)[(m.p)] == 58 {
			goto tr25
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
	_testEof125:
		m.cs = 125
		goto _testEof
	_testEof9:
		m.cs = 9
		goto _testEof
	_testEof126:
		m.cs = 126
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
	_testEof33:
		m.cs = 33
		goto _testEof
	_testEof130:
		m.cs = 130
		goto _testEof
	_testEof131:
		m.cs = 131
		goto _testEof
	_testEof34:
		m.cs = 34
		goto _testEof
	_testEof132:
		m.cs = 132
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
	_testEof133:
		m.cs = 133
		goto _testEof
	_testEof44:
		m.cs = 44
		goto _testEof
	_testEof134:
		m.cs = 134
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
	_testEof61:
		m.cs = 61
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
	_testEof70:
		m.cs = 70
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
	_testEof135:
		m.cs = 135
		goto _testEof
	_testEof85:
		m.cs = 85
		goto _testEof
	_testEof136:
		m.cs = 136
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
	_testEof101:
		m.cs = 101
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
	_testEof106:
		m.cs = 106
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
	_testEof110:
		m.cs = 110
		goto _testEof
	_testEof111:
		m.cs = 111
		goto _testEof
	_testEof112:
		m.cs = 112
		goto _testEof
	_testEof113:
		m.cs = 113
		goto _testEof
	_testEof114:
		m.cs = 114
		goto _testEof
	_testEof115:
		m.cs = 115
		goto _testEof
	_testEof117:
		m.cs = 117
		goto _testEof
	_testEof118:
		m.cs = 118
		goto _testEof
	_testEof119:
		m.cs = 119
		goto _testEof
	_testEof120:
		m.cs = 120
		goto _testEof
	_testEof137:
		m.cs = 137
		goto _testEof
	_testEof121:
		m.cs = 121
		goto _testEof
	_testEof138:
		m.cs = 138
		goto _testEof
	_testEof122:
		m.cs = 122
		goto _testEof
	_testEof123:
		m.cs = 123
		goto _testEof
	_testEof124:
		m.cs = 124
		goto _testEof
	_testEof127:
		m.cs = 127
		goto _testEof
	_testEof14:
		m.cs = 14
		goto _testEof
	_testEof15:
		m.cs = 15
		goto _testEof
	_testEof128:
		m.cs = 128
		goto _testEof
	_testEof16:
		m.cs = 16
		goto _testEof
	_testEof17:
		m.cs = 17
		goto _testEof
	_testEof129:
		m.cs = 129
		goto _testEof
	_testEof18:
		m.cs = 18
		goto _testEof
	_testEof19:
		m.cs = 19
		goto _testEof
	_testEof20:
		m.cs = 20
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
	_testEof29:
		m.cs = 29
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

	_testEof:
		{
		}
		if (m.p) == (m.eof) {
			switch m.cs {
			case 2, 3, 4, 13, 36, 37, 38, 39, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71, 72, 73, 74, 75, 77, 78, 79, 80, 89, 90, 91, 92, 93, 94, 95, 96, 97, 98, 99, 100, 101, 102, 103, 104, 105, 106, 107, 108, 109, 110, 111, 112, 113, 114, 115:

				if m.pe > 0 {
					if m.p != m.pe {
						m.err = m.emitErrorOnCurrentCharacter(ErrType)
					} else {
						// assert(m.p == m.pe)
						m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
					}
				}

			case 5, 6, 12, 40, 41, 47, 81, 82, 88, 117, 118, 124:

				if m.err == nil {
					m.err = m.emitErrorOnCurrentCharacter(ErrColon)
				}

			case 7, 42, 83, 119:

				if m.err == nil {
					m.err = m.emitErrorOnCurrentCharacter(ErrDescriptionInit)
				}

			case 8, 43, 84, 120:

				if m.p < m.pe && m.data[m.p] == 10 {
					m.err = m.emitError(ErrNewline, m.p+1)
				} else {
					// assert(m.p == m.pe)
					m.err = m.emitErrorOnPreviousCharacter(ErrDescription)
				}

			case 9, 44, 85, 121:

				m.err = m.emitErrorWithoutCharacter(ErrMissingBlankLineAtBeginning)

			case 125, 133, 135, 137:

				output.descr = string(m.text())
				m.emitInfo("valid commit message description", "description", output.descr)

			case 130:

				output.footers[m.currentFooterKey] = append(output.footers[m.currentFooterKey], string(m.text()))
				m.emitInfo("valid commit message footer trailer", m.currentFooterKey, string(m.text()))

			case 132:

				// Append newlines
				for m.countNewlines > 0 {
					output.body += "\n"
					m.countNewlines--
					m.emitInfo("valid commit message body content", "body", "\n")
				}
				// Append body content
				output.body += string(m.text())
				m.emitInfo("valid commit message body content", "body", string(m.text()))

			case 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32:

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
						goto st34
					}
				} else {
					// A rewind happens when an error while parsing a footer trailer is encountered
					// If this is not the first footer trailer the parser can't go back to parse body content again
					// Thus, emit an error
					if m.p != m.pe {
						m.err = m.emitErrorOnCurrentCharacter(ErrTrailer)
					} else {
						// assert(m.p == m.pe)
						m.err = m.emitErrorOnPreviousCharacter(ErrTrailerIncomplete)
					}
				}

			case 1, 35, 76, 116:

				m.err = m.emitErrorWithoutCharacter(ErrEmpty)

				if m.pe > 0 {
					if m.p != m.pe {
						m.err = m.emitErrorOnCurrentCharacter(ErrType)
					} else {
						// assert(m.p == m.pe)
						m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
					}
				}

			case 10, 11, 45, 46, 86, 87, 122, 123:

				if m.p < m.pe {
					m.err = m.emitErrorOnCurrentCharacter(ErrScope)
				}

				// assert(m.p == m.pe)
				m.err = m.emitErrorOnPreviousCharacter(ErrScopeIncomplete)

			case 34:

				// Append newlines
				for m.countNewlines > 0 {
					output.body += "\n"
					m.countNewlines--
					m.emitInfo("valid commit message body content", "body", "\n")
				}
				// Append body content
				output.body += string(m.text())
				m.emitInfo("valid commit message body content", "body", string(m.text()))

				m.emitDebug("try to parse a footer trailer token", "pos", m.p)
				{
					goto st127
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
