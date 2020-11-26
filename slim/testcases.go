package slim

import (
	"fmt"

	"github.com/leodido/go-conventionalcommits"
	cctesting "github.com/leodido/go-conventionalcommits/testing"
)

type testCase struct {
	title        string
	input        []byte
	ok           bool
	value        conventionalcommits.Message
	partialValue conventionalcommits.Message
	errorString  string
}

var testCases = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("f"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "f", 1),
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("fx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("fit"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("fea"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "a", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("feax"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("feat"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("fi:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("fea:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: x"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "x",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "x",
			},
		},
		"",
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("fix>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("feat?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("fix(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 9),
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("fix(aaa): bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		"",
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
			},
		},
		"",
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("f description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("fix(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 10),
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("fix(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 11),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("fix(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("fix(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 14),
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("fix(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 4),
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("fix(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 5),
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("feat(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 8),
	},
}

var testCasesForFalcoTypes = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("c"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 1),
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("bx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-revert",
		[]byte("re"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("net"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("bui"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "i", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("docx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("perz"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "z", 3),
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("test"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("ch:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("upd:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: w"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "w",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "w",
			},
		},
		"",
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message-rule",
		[]byte("rule: super secure rule"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "rule",
				Description: "super secure rule",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "rule",
				Description: "super secure rule",
			},
		},
		"",
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("new>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("perf?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-5-chars",
		[]byte("build?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 5),
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("new(xyz): ccc"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "new",
				Scope:       cctesting.StringAddress("xyz"),
				Description: "ccc",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "new",
				Scope:       cctesting.StringAddress("xyz"),
				Description: "ccc",
			},
		},
		"",
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
			},
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope-feat",
		[]byte("feat(aaa)!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "feat",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "feat",
				Scope:       cctesting.StringAddress("aaa"),
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
			},
		},
		"",
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		&ConventionalCommit{
			Minimal: conventionalcommits.Minimal{
				Type:        "fix",
				Description: "bbb",
				Exclamation: true,
			},
		},
		"",
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("new(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 9),
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("c description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("new(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 10),
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("new(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 11),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("revert(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 16),
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("ci(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("chore(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 6),
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("perf(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrMalformedScope+ColumnPositionTemplate, "(", 6),
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("new(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 7),
	},
}
