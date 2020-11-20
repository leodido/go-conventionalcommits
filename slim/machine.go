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
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
	// ErrEarly ..
	ErrEarly = "early exit after '%s' character"
)

const start int = 1
const firstFinal int = 76

const enMinimalTypes int = 2
const enConventionalTypes int = 7
const enFalcoTypes int = 40
const enScope int = 80
const enBreaking int = 82
const enSeparator int = 74
const enDescription int = 75
const enFail int = 87
const enMain int = 1

type machine struct {
	data       []byte
	cs         int
	p, pe, eof int
	pb         int
	err        error
	bestEffort bool
	typeConfig conventionalcommits.TypeConfig
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

func (m *machine) emitErrorWithoutCharacter(messageTemplate string) error {
	return fmt.Errorf(messageTemplate+ColumnPositionTemplate, m.p)
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
		case 76:
			goto st76
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
		case 77:
			goto st77
		case 6:
			goto st6
		case 7:
			goto st7
		case 8:
			goto st8
		case 9:
			goto st9
		case 10:
			goto st10
		case 11:
			goto st11
		case 78:
			goto st78
		case 12:
			goto st12
		case 13:
			goto st13
		case 14:
			goto st14
		case 15:
			goto st15
		case 16:
			goto st16
		case 17:
			goto st17
		case 18:
			goto st18
		case 19:
			goto st19
		case 20:
			goto st20
		case 21:
			goto st21
		case 22:
			goto st22
		case 23:
			goto st23
		case 24:
			goto st24
		case 25:
			goto st25
		case 26:
			goto st26
		case 27:
			goto st27
		case 28:
			goto st28
		case 29:
			goto st29
		case 30:
			goto st30
		case 31:
			goto st31
		case 32:
			goto st32
		case 33:
			goto st33
		case 34:
			goto st34
		case 35:
			goto st35
		case 36:
			goto st36
		case 37:
			goto st37
		case 38:
			goto st38
		case 39:
			goto st39
		case 40:
			goto st40
		case 41:
			goto st41
		case 42:
			goto st42
		case 43:
			goto st43
		case 44:
			goto st44
		case 79:
			goto st79
		case 45:
			goto st45
		case 46:
			goto st46
		case 47:
			goto st47
		case 48:
			goto st48
		case 49:
			goto st49
		case 50:
			goto st50
		case 51:
			goto st51
		case 52:
			goto st52
		case 53:
			goto st53
		case 54:
			goto st54
		case 55:
			goto st55
		case 56:
			goto st56
		case 57:
			goto st57
		case 58:
			goto st58
		case 59:
			goto st59
		case 60:
			goto st60
		case 61:
			goto st61
		case 62:
			goto st62
		case 63:
			goto st63
		case 64:
			goto st64
		case 65:
			goto st65
		case 66:
			goto st66
		case 67:
			goto st67
		case 68:
			goto st68
		case 69:
			goto st69
		case 70:
			goto st70
		case 71:
			goto st71
		case 74:
			goto st74
		case 84:
			goto st84
		case 75:
			goto st75
		case 85:
			goto st85
		case 86:
			goto st86
		case 80:
			goto st80
		case 72:
			goto st72
		case 73:
			goto st73
		case 81:
			goto st81
		case 82:
			goto st82
		case 83:
			goto st83
		case 87:
			goto st87
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof
		}
	_resume:
		switch m.cs {
		case 1:
			goto stCase1
		case 76:
			goto stCase76
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
		case 77:
			goto stCase77
		case 6:
			goto stCase6
		case 7:
			goto stCase7
		case 8:
			goto stCase8
		case 9:
			goto stCase9
		case 10:
			goto stCase10
		case 11:
			goto stCase11
		case 78:
			goto stCase78
		case 12:
			goto stCase12
		case 13:
			goto stCase13
		case 14:
			goto stCase14
		case 15:
			goto stCase15
		case 16:
			goto stCase16
		case 17:
			goto stCase17
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
		case 79:
			goto stCase79
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
		case 74:
			goto stCase74
		case 84:
			goto stCase84
		case 75:
			goto stCase75
		case 85:
			goto stCase85
		case 86:
			goto stCase86
		case 80:
			goto stCase80
		case 72:
			goto stCase72
		case 73:
			goto stCase73
		case 81:
			goto stCase81
		case 82:
			goto stCase82
		case 83:
			goto stCase83
		case 87:
			goto stCase87
		}
		goto stOut
	st1:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof1
		}
	stCase1:
		goto tr0
	tr0:
		m.cs = 76

		(m.p)--

		switch m.typeConfig {
		case conventionalcommits.TypesMinimal:
			m.cs = 2
			break
		case conventionalcommits.TypesConventional:
			m.cs = 7
			break
		case conventionalcommits.TypesFalco:
			m.cs = 40
			break
		}

		goto _again
	st76:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof76
		}
	stCase76:
		goto st0
	tr3:

		if m.p != m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrType)
		} else {
			m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
		}

		goto st0
	tr77:

		m.err = m.emitErrorOnCurrentCharacter(ErrColon)

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
			goto st77
		}
		goto tr3
	st77:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
			{
				goto st87
			}
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof77
		}
	stCase77:

		output._type = string(m.text())

		fmt.Println("goto scope")
		(m.p)--

		{
			goto st80
		}

		goto st0
	st6:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof6
		}
	stCase6:
		if (m.data)[(m.p)] == 120 {
			goto st77
		}
		goto tr3
	st7:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof7
		}
	stCase7:
		switch (m.data)[(m.p)] {
		case 98:
			goto tr8
		case 99:
			goto tr9
		case 100:
			goto tr10
		case 102:
			goto tr11
		case 112:
			goto tr12
		case 114:
			goto tr13
		case 115:
			goto tr14
		case 116:
			goto tr15
		}
		goto st0
	tr8:

		m.pb = m.p

		goto st8
	st8:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof8
		}
	stCase8:
		if (m.data)[(m.p)] == 117 {
			goto st9
		}
		goto tr3
	st9:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof9
		}
	stCase9:
		if (m.data)[(m.p)] == 105 {
			goto st10
		}
		goto tr3
	st10:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof10
		}
	stCase10:
		if (m.data)[(m.p)] == 108 {
			goto st11
		}
		goto tr3
	st11:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof11
		}
	stCase11:
		if (m.data)[(m.p)] == 100 {
			goto st78
		}
		goto tr3
	st78:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
			{
				goto st87
			}
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof78
		}
	stCase78:

		output._type = string(m.text())

		fmt.Println("goto scope")
		(m.p)--

		{
			goto st80
		}

		goto st0
	tr9:

		m.pb = m.p

		goto st12
	st12:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof12
		}
	stCase12:
		switch (m.data)[(m.p)] {
		case 104:
			goto st13
		case 105:
			goto st78
		}
		goto tr3
	st13:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof13
		}
	stCase13:
		if (m.data)[(m.p)] == 111 {
			goto st14
		}
		goto tr3
	st14:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof14
		}
	stCase14:
		if (m.data)[(m.p)] == 114 {
			goto st15
		}
		goto tr3
	st15:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof15
		}
	stCase15:
		if (m.data)[(m.p)] == 101 {
			goto st78
		}
		goto tr3
	tr10:

		m.pb = m.p

		goto st16
	st16:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof16
		}
	stCase16:
		if (m.data)[(m.p)] == 111 {
			goto st17
		}
		goto tr3
	st17:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof17
		}
	stCase17:
		if (m.data)[(m.p)] == 99 {
			goto st18
		}
		goto tr3
	st18:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof18
		}
	stCase18:
		if (m.data)[(m.p)] == 115 {
			goto st78
		}
		goto tr3
	tr11:

		m.pb = m.p

		goto st19
	st19:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof19
		}
	stCase19:
		switch (m.data)[(m.p)] {
		case 101:
			goto st20
		case 105:
			goto st22
		}
		goto tr3
	st20:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof20
		}
	stCase20:
		if (m.data)[(m.p)] == 97 {
			goto st21
		}
		goto tr3
	st21:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof21
		}
	stCase21:
		if (m.data)[(m.p)] == 116 {
			goto st78
		}
		goto tr3
	st22:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof22
		}
	stCase22:
		if (m.data)[(m.p)] == 120 {
			goto st78
		}
		goto tr3
	tr12:

		m.pb = m.p

		goto st23
	st23:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof23
		}
	stCase23:
		if (m.data)[(m.p)] == 101 {
			goto st24
		}
		goto tr3
	st24:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof24
		}
	stCase24:
		if (m.data)[(m.p)] == 114 {
			goto st25
		}
		goto tr3
	st25:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof25
		}
	stCase25:
		if (m.data)[(m.p)] == 102 {
			goto st78
		}
		goto tr3
	tr13:

		m.pb = m.p

		goto st26
	st26:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof26
		}
	stCase26:
		if (m.data)[(m.p)] == 101 {
			goto st27
		}
		goto tr3
	st27:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof27
		}
	stCase27:
		switch (m.data)[(m.p)] {
		case 102:
			goto st28
		case 118:
			goto st33
		}
		goto tr3
	st28:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof28
		}
	stCase28:
		if (m.data)[(m.p)] == 97 {
			goto st29
		}
		goto tr3
	st29:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof29
		}
	stCase29:
		if (m.data)[(m.p)] == 99 {
			goto st30
		}
		goto tr3
	st30:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof30
		}
	stCase30:
		if (m.data)[(m.p)] == 116 {
			goto st31
		}
		goto tr3
	st31:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof31
		}
	stCase31:
		if (m.data)[(m.p)] == 111 {
			goto st32
		}
		goto tr3
	st32:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof32
		}
	stCase32:
		if (m.data)[(m.p)] == 114 {
			goto st78
		}
		goto tr3
	st33:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof33
		}
	stCase33:
		if (m.data)[(m.p)] == 101 {
			goto st34
		}
		goto tr3
	st34:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof34
		}
	stCase34:
		if (m.data)[(m.p)] == 114 {
			goto st21
		}
		goto tr3
	tr14:

		m.pb = m.p

		goto st35
	st35:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof35
		}
	stCase35:
		if (m.data)[(m.p)] == 116 {
			goto st36
		}
		goto tr3
	st36:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof36
		}
	stCase36:
		if (m.data)[(m.p)] == 121 {
			goto st37
		}
		goto tr3
	st37:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof37
		}
	stCase37:
		if (m.data)[(m.p)] == 108 {
			goto st15
		}
		goto tr3
	tr15:

		m.pb = m.p

		goto st38
	st38:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof38
		}
	stCase38:
		if (m.data)[(m.p)] == 101 {
			goto st39
		}
		goto tr3
	st39:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof39
		}
	stCase39:
		if (m.data)[(m.p)] == 115 {
			goto st21
		}
		goto tr3
	st40:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof40
		}
	stCase40:
		switch (m.data)[(m.p)] {
		case 98:
			goto tr41
		case 99:
			goto tr42
		case 100:
			goto tr43
		case 102:
			goto tr44
		case 110:
			goto tr45
		case 112:
			goto tr46
		case 114:
			goto tr47
		case 116:
			goto tr48
		case 117:
			goto tr49
		}
		goto st0
	tr41:

		m.pb = m.p

		goto st41
	st41:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof41
		}
	stCase41:
		if (m.data)[(m.p)] == 117 {
			goto st42
		}
		goto tr3
	st42:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof42
		}
	stCase42:
		if (m.data)[(m.p)] == 105 {
			goto st43
		}
		goto tr3
	st43:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof43
		}
	stCase43:
		if (m.data)[(m.p)] == 108 {
			goto st44
		}
		goto tr3
	st44:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof44
		}
	stCase44:
		if (m.data)[(m.p)] == 100 {
			goto st79
		}
		goto tr3
	st79:

		if (m.p + 1) == m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
			{
				goto st87
			}
		}

		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof79
		}
	stCase79:

		output._type = string(m.text())

		fmt.Println("goto scope")
		(m.p)--

		{
			goto st80
		}

		goto st0
	tr42:

		m.pb = m.p

		goto st45
	st45:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof45
		}
	stCase45:
		switch (m.data)[(m.p)] {
		case 104:
			goto st46
		case 105:
			goto st79
		}
		goto tr3
	st46:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof46
		}
	stCase46:
		if (m.data)[(m.p)] == 111 {
			goto st47
		}
		goto tr3
	st47:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof47
		}
	stCase47:
		if (m.data)[(m.p)] == 114 {
			goto st48
		}
		goto tr3
	st48:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof48
		}
	stCase48:
		if (m.data)[(m.p)] == 101 {
			goto st79
		}
		goto tr3
	tr43:

		m.pb = m.p

		goto st49
	st49:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof49
		}
	stCase49:
		if (m.data)[(m.p)] == 111 {
			goto st50
		}
		goto tr3
	st50:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof50
		}
	stCase50:
		if (m.data)[(m.p)] == 99 {
			goto st51
		}
		goto tr3
	st51:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof51
		}
	stCase51:
		if (m.data)[(m.p)] == 115 {
			goto st79
		}
		goto tr3
	tr44:

		m.pb = m.p

		goto st52
	st52:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof52
		}
	stCase52:
		switch (m.data)[(m.p)] {
		case 101:
			goto st53
		case 105:
			goto st55
		}
		goto tr3
	st53:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof53
		}
	stCase53:
		if (m.data)[(m.p)] == 97 {
			goto st54
		}
		goto tr3
	st54:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof54
		}
	stCase54:
		if (m.data)[(m.p)] == 116 {
			goto st79
		}
		goto tr3
	st55:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof55
		}
	stCase55:
		if (m.data)[(m.p)] == 120 {
			goto st79
		}
		goto tr3
	tr45:

		m.pb = m.p

		goto st56
	st56:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof56
		}
	stCase56:
		if (m.data)[(m.p)] == 101 {
			goto st57
		}
		goto tr3
	st57:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof57
		}
	stCase57:
		if (m.data)[(m.p)] == 119 {
			goto st79
		}
		goto tr3
	tr46:

		m.pb = m.p

		goto st58
	st58:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof58
		}
	stCase58:
		if (m.data)[(m.p)] == 101 {
			goto st59
		}
		goto tr3
	st59:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof59
		}
	stCase59:
		if (m.data)[(m.p)] == 114 {
			goto st60
		}
		goto tr3
	st60:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof60
		}
	stCase60:
		if (m.data)[(m.p)] == 102 {
			goto st79
		}
		goto tr3
	tr47:

		m.pb = m.p

		goto st61
	st61:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof61
		}
	stCase61:
		switch (m.data)[(m.p)] {
		case 101:
			goto st62
		case 117:
			goto st65
		}
		goto tr3
	st62:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof62
		}
	stCase62:
		if (m.data)[(m.p)] == 118 {
			goto st63
		}
		goto tr3
	st63:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof63
		}
	stCase63:
		if (m.data)[(m.p)] == 101 {
			goto st64
		}
		goto tr3
	st64:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof64
		}
	stCase64:
		if (m.data)[(m.p)] == 114 {
			goto st54
		}
		goto tr3
	st65:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof65
		}
	stCase65:
		if (m.data)[(m.p)] == 108 {
			goto st48
		}
		goto tr3
	tr48:

		m.pb = m.p

		goto st66
	st66:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof66
		}
	stCase66:
		if (m.data)[(m.p)] == 101 {
			goto st67
		}
		goto tr3
	st67:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof67
		}
	stCase67:
		if (m.data)[(m.p)] == 115 {
			goto st54
		}
		goto tr3
	tr49:

		m.pb = m.p

		goto st68
	st68:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof68
		}
	stCase68:
		if (m.data)[(m.p)] == 112 {
			goto st69
		}
		goto tr3
	st69:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof69
		}
	stCase69:
		if (m.data)[(m.p)] == 100 {
			goto st70
		}
		goto tr3
	st70:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof70
		}
	stCase70:
		if (m.data)[(m.p)] == 97 {
			goto st71
		}
		goto tr3
	st71:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof71
		}
	stCase71:
		if (m.data)[(m.p)] == 116 {
			goto st48
		}
		goto tr3
	st74:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof74
		}
	stCase74:
		if (m.data)[(m.p)] == 58 {
			goto st84
		}
		goto tr77
	st84:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof84
		}
	stCase84:

		fmt.Println("goto description")
		(m.p)--

		{
			goto st75
		}

		goto st0
	st75:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof75
		}
	stCase75:
		if (m.data)[(m.p)] == 32 {
			goto st85
		}
		goto st0
	tr83:

		m.pb = m.p

		goto st85
	st85:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof85
		}
	stCase85:
		if (m.data)[(m.p)] == 32 {
			goto tr83
		}
		goto tr82
	tr82:

		m.pb = m.p

		goto st86
	st86:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof86
		}
	stCase86:
		goto st86
	st80:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof80
		}
	stCase80:

		fmt.Println("goto breaking")
		(m.p)--

		{
			goto st82
		}

		if (m.data)[(m.p)] == 40 {
			goto st72
		}
		goto st0
	st72:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof72
		}
	stCase72:
		if (m.data)[(m.p)] == 41 {
			goto tr74
		}
		goto tr73
	tr73:

		m.pb = m.p

		goto st73
	st73:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof73
		}
	stCase73:
		if (m.data)[(m.p)] == 41 {
			goto tr76
		}
		goto st73
	tr74:

		m.pb = m.p

		output.scope = string(m.text())

		goto st81
	tr76:

		output.scope = string(m.text())

		goto st81
	st81:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof81
		}
	stCase81:

		fmt.Println("goto breaking")
		(m.p)--

		{
			goto st82
		}

		if (m.data)[(m.p)] == 41 {
			goto tr76
		}
		goto st73
	st82:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof82
		}
	stCase82:

		fmt.Println("goto separator")
		(m.p)--

		{
			goto st74
		}

		if (m.data)[(m.p)] == 33 {
			goto tr81
		}
		goto st0
	tr81:

		m.pb = m.p

		goto st83
	st83:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof83
		}
	stCase83:

		fmt.Println("goto separator")
		(m.p)--

		{
			goto st74
		}

		goto st0
	st87:
		if (m.p)++; (m.p) == (m.pe) {
			goto _testEof87
		}
	stCase87:
		switch (m.data)[(m.p)] {
		case 10:
			goto st0
		case 13:
			goto st0
		}
		goto st87
	stOut:
	_testEof1:
		m.cs = 1
		goto _testEof
	_testEof76:
		m.cs = 76
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
	_testEof77:
		m.cs = 77
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
	_testEof9:
		m.cs = 9
		goto _testEof
	_testEof10:
		m.cs = 10
		goto _testEof
	_testEof11:
		m.cs = 11
		goto _testEof
	_testEof78:
		m.cs = 78
		goto _testEof
	_testEof12:
		m.cs = 12
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
	_testEof16:
		m.cs = 16
		goto _testEof
	_testEof17:
		m.cs = 17
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
	_testEof79:
		m.cs = 79
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
	_testEof74:
		m.cs = 74
		goto _testEof
	_testEof84:
		m.cs = 84
		goto _testEof
	_testEof75:
		m.cs = 75
		goto _testEof
	_testEof85:
		m.cs = 85
		goto _testEof
	_testEof86:
		m.cs = 86
		goto _testEof
	_testEof80:
		m.cs = 80
		goto _testEof
	_testEof72:
		m.cs = 72
		goto _testEof
	_testEof73:
		m.cs = 73
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
	_testEof87:
		m.cs = 87
		goto _testEof

	_testEof:
		{
		}
		if (m.p) == (m.eof) {
			switch m.cs {
			case 1:

				m.err = m.emitErrorWithoutCharacter(ErrEmpty)

			case 3, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65, 66, 67, 68, 69, 70, 71:

				if m.p != m.pe {
					m.err = m.emitErrorOnCurrentCharacter(ErrType)
				} else {
					m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
				}

			case 74:

				m.err = m.emitErrorOnCurrentCharacter(ErrColon)

			case 77, 78, 79:

				output._type = string(m.text())

			case 86:

				output.descr = string(m.text())

			case 83:

				output.exclamation = true

			case 85:

				m.pb = m.p

				fmt.Println("FINE")

				output.descr = string(m.text())

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

// WithTypes ...
func (m *machine) WithTypes(t conventionalcommits.TypeConfig) {
	m.typeConfig = t
}
