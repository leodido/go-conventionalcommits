// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

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
	bump         *conventionalcommits.VersionBump
}

var UnknownVersion = conventionalcommits.UnknownVersion
var PatchVersion = conventionalcommits.PatchVersion
var MinorVersion = conventionalcommits.MinorVersion
var MajorVersion = conventionalcommits.MajorVersion

var testCases = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
		nil,
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("f"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "f", 1),
		nil,
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("fx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
		nil,
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
		nil,
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("fit"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
		nil,
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
		nil,
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("fea"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "a", 3),
		nil,
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("feax"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
		nil,
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
		nil,
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("feat"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
		nil,
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("fi:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
		nil,
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("fea:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
		nil,
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: x"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("fix>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("feat?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
		nil,
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("fix(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 9),
		nil,
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("fix(aaa): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
		nil,
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
		nil,
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
		nil,
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
		nil,
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
		nil,
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("f description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
		nil,
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("fix(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 10),
		nil,
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("fix(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 11),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("fix(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("fix(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 14),
		nil,
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("fix(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 4),
		nil,
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("fix(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 5),
		nil,
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
		nil,
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("feat(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 8),
		nil,
	},
	// INVALID / newline rather than whitespace in description
	{
		"invalid-newline-rather-than-whitespace-description",
		[]byte("feat(az):\x0A description on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "\n", 9),
		nil,
	},
	// INVALID / newline after whitespace in description
	{
		"invalid-newline-after-whitespace-description",
		[]byte("feat(az): \x0Adescription on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrNewline+ColumnPositionTemplate, 11),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description",
		[]byte("feat(az): new\x0Aline"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Description: "new",
			TypeConfig:  0,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 14),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description-2",
		[]byte("feat(az)!: bla\x0Al"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  0,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"description-ending-with-single-newline",
		[]byte("feat(az)!: bla\x0A"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  0,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
		nil,
	},
	// VALID / multi-line body is valid (after a blank line)
	{
		"valid-with-multi-line-body",
		[]byte(`fix: x

see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / multi-line body ending with multiple blank lines (they gets discarded) is valid
	{
		"valid-with-multi-line-body-ending-extras-blank-lines",
		[]byte(`fix: x

see the issue for details

on typos fixed.

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / multi-line body starting with many extra blank lines is valid
	{
		"valid-with-multi-line-body-after-two-extra-blank-lines",
		[]byte(`fix: magic



see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / multi-line body starting and ending with many extra blank lines is valid
	{
		"valid-with-multi-line-body-with-extra-blank-lines-before-and-after",
		[]byte(`fix: magic



see the issue for details

on typos fixed.


`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "magic",
			Body:        cctesting.StringAddress("\n\nsee the issue for details\n\non typos fixed."),
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / single line body (after blank line) is valid
	{
		"valid-with-single-line-body",
		[]byte(`fix: correct minor typos in code

see the issue for details.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / empty body is okay (it's optional)
	{
		"valid-with-empty-body",
		[]byte(`fix: correct something

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / multiple blank lines body is okay (it's considered empty)
	{
		"valid-with-multiple-blank-lines-body",
		[]byte(`fix: descr





`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "descr",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "descr",
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / only footer
	{
		"valid-with-footer-only",
		[]byte(`fix: only footer

Fixes #3
Signed-off-by: Leo`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / only footer after many blank lines (that gets ignored)
	{
		"valid-with-footer-only-after-many-blank-lines",
		[]byte(`fix: only footer




Fixes #3
Signed-off-by: Leo`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / only footer ending with many blank lines (that gets ignored)
	{
		"valid-with-footer-only-ending-with-many-blank-lines",
		[]byte(`fix: only footer

Fixes #3
Signed-off-by: Leo


`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes":         {"3"},
				"signed-off-by": {"Leo"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / only footer containing repetitions
	{
		"valid-with-footer-containing-repetitions",
		[]byte(`fix: only footer

Fixes #3
Fixes #4
Fixes #5`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes": {"3", "4", "5"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "only footer",
			Footers: map[string][]string{
				"fixes": {"3", "4", "5"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / Multi-line body with extras blank lines after and footer with multiple trailers
	{
		"valid-with-multi-line-body-containing-extra-blank-lines-inside-and-after-plus-footer-many-trailers",
		[]byte(`fix: sarah

FUCK

COVID-19.
This is the only message I have in my mind

right now.



Fixes #22
Co-authored-by: My other personality <persona@email.com>
Signed-off-by: Leonardo Di Donato <some@email.com>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\n\nCOVID-19.\nThis is the only message I have in my mind\n\nright now."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\n\nCOVID-19.\nThis is the only message I have in my mind\n\nright now."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / Multi-line body with newlines inside and many blank lines after and footer with multiple trailers
	{
		"valid-with-multi-line-body-and-extra-blank-lines-after-plus-footer-many-trailers",
		[]byte(`fix: sarah

FUCK
COVID-19.
This is the only message I have in my mind
right
now.



Fixes #22
Co-authored-by: My other personality <persona@email.com>
Signed-off-by: Leonardo Di Donato <some@email.com>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			Body:        cctesting.StringAddress("FUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / Multi-line body with newlines inside and many blank lines before it, plus footer with multiple trailers
	{
		"valid-with-multi-line-body-and-extra-blank-lines-before-plus-footer-many-trailers",
		[]byte(`fix: sarah



FUCK
COVID-19.
This is the only message I have in my mind
right
now.



Fixes #22
Co-authored-by: My other personality <persona@email.com>
Signed-off-by: Leonardo Di Donato <some@email.com>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			// First blank line ("\n\n") gets ignored
			Body: cctesting.StringAddress("\n\nFUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "sarah",
			// First blank line ("\n\n") gets ignored
			Body: cctesting.StringAddress("\n\nFUCK\nCOVID-19.\nThis is the only message I have in my mind\nright\nnow."),
			Footers: map[string][]string{
				"fixes":          {"22"},
				"co-authored-by": {"My other personality <persona@email.com>"},
				"signed-off-by":  {"Leonardo Di Donato <some@email.com>"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
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
		nil,
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("c"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 1),
		nil,
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("bx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
		nil,
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
		nil,
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-revert",
		[]byte("re"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
		nil,
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("net"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
		nil,
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
		nil,
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("bui"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "i", 3),
		nil,
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("docx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
		nil,
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("perz"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "z", 3),
		nil,
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
		nil,
	},
	// INVALID / missing colon after type feat
	{
		"invalid-after-valid-type-feat",
		[]byte("test"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
		nil,
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("ch:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
		nil,
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("upd:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
		nil,
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / minimal commit message with minor bump for new
	{
		"valid-minimal-commit-message-with-minor-bump-for-new",
		[]byte("new: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Description: "w",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Description: "w",
			TypeConfig:  2,
		},
		"",
		&PatchVersion,
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message-rule",
		[]byte("rule: super secure rule"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "rule",
			Description: "super secure rule",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "rule",
			Description: "super secure rule",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / minimal commit message with uppercase type
	{
		"valid-minimal-commit-message-rule-with-uppercase-type",
		[]byte("RULE: super secure rule"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "rule",
			Description: "super secure rule",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "rule",
			Description: "super secure rule",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("new>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
		nil,
	},
	// INVALID / missing colon after valid uppercase commit message type
	{
		"missing-colon-after-uppercase-type-3-chars",
		[]byte("NEW>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("perf?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
		nil,
	},
	// INVALID / missing colon after valid uppercase commit message type
	{
		"missing-colon-after-uppercase-type-4-chars",
		[]byte("PERF?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-5-chars",
		[]byte("build?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 5),
		nil,
	},
	// INVALID / missing colon after valid uppercase commit message type
	{
		"missing-colon-after-uppercase-type-5-chars",
		[]byte("BUILD?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 5),
		nil,
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("new(xyz): ccc"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / uppercase type + scope + description
	{
		"valid-with-uppercase-type-and-scope",
		[]byte("NEW(xyz): ccc"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope-feat",
		[]byte("feat(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / uppercase type + scope + breaking + description
	{
		"valid-breaking-with-scope-feat-uppercase-type",
		[]byte("FEAT(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / empty scope is ignored (uppercase type)
	{
		"valid-empty-scope-is-ignored-uppercase-type",
		[]byte("FIX(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
		nil,
	},
	// INVALID / missing whitespace after colon (with breaking, uppercase type)
	{
		"invalid-missing-ws-after-colon-with-breaking-uppercase-type",
		[]byte("FIX!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
		nil,
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
		nil,
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
		nil,
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
		nil,
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("new(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 9),
		nil,
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
		nil,
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("c description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
		nil,
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("new(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 10),
		nil,
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("new(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 11),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("revert(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 16),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("ci(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
		nil,
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("chore(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 6),
		nil,
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("perf(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 6),
		nil,
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
		nil,
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("new(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 7),
		nil,
	},
	// INVALID / newline rather than whitespace in description
	{
		"invalid-newline-rather-than-whitespace-description",
		[]byte("perf(ax):\x0A description on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "\n", 9),
		nil,
	},
	// INVALID / newline after whitespace in description
	{
		"invalid-newline-after-whitespace-description",
		[]byte("feat(az): \x0Adescription on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrNewline+ColumnPositionTemplate, 11),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description",
		[]byte("feat(ae): new\x0Aline"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("ae"),
			Description: "new",
			TypeConfig:  2,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 14),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"invalid-newline-in-description-2",
		[]byte("docs(az)!: bla\x0Al"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "docs",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  2,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"description-ending-with-single-newline",
		[]byte("docs(az)!: bla\x0A"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "docs",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  2,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
		nil,
	},
	// VALID
	{
		"valid-with-multiline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID
	{
		"valid-with-singleline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID
	{
		"valid-with-empty-body",
		[]byte(`fix: correct something

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  2,
		},
		"",
		nil,
	},
	// VALID
	{
		"valid-with-multiple-blank-lines-body",
		[]byte(`fix: correct something



`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  2,
		},
		"",
		nil,
	},
}

var testCasesForConventionalTypes = []testCase{
	// INVALID / empty
	{
		"empty",
		[]byte(""),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEmpty+ColumnPositionTemplate, 0),
		nil,
	},
	// INVALID / invalid type (1 char)
	{
		"invalid-type-1-char",
		[]byte("c"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 1),
		nil,
	},
	// INVALID / invalid type (2 char)
	{
		"invalid-type-2-char",
		[]byte("bx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 1),
		nil,
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-feat",
		[]byte("fe"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
		nil,
	},
	// INVALID / invalid type (2 char) with almost valid type
	{
		"invalid-type-2-char-revert",
		[]byte("re"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "e", 2),
		nil,
	},
	// INVALID / invalid type (3 char)
	{
		"invalid-type-3-char",
		[]byte("fit"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "t", 2),
		nil,
	},
	// INVALID / invalid type (3 char) again
	{
		"invalid-type-3-char-feat",
		[]byte("fei"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "i", 2),
		nil,
	},
	// INVALID / invalid type (3 char) with almost valid type
	{
		"invalid-type-3-char-feat",
		[]byte("bui"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "i", 3),
		nil,
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("tesx"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "x", 3),
		nil,
	},
	// INVALID / invalid type (4 char) with almost valid type
	{
		"invalid-type-4-char-refactor",
		[]byte("refa"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "a", 4),
		nil,
	},
	// INVALID / invalid type (4 char)
	{
		"invalid-type-4-char",
		[]byte("perz"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "z", 3),
		nil,
	},
	// INVALID / invalid type (5 char) with almost valid type
	{
		"invalid-type-5-char-refactor",
		[]byte("refac"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "c", 5),
		nil,
	},
	// INVALID / invalid type (6 char) with almost valid type
	{
		"invalid-type-6-char-refactor",
		[]byte("refact"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "t", 6),
		nil,
	},
	// INVALID / invalid type (7 char) with almost valid type
	{
		"invalid-type-7-char-refactor",
		[]byte("refacto"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrTypeIncomplete+ColumnPositionTemplate, "o", 7),
		nil,
	},
	// INVALID / missing colon after type fix
	{
		"invalid-after-valid-type-fix",
		[]byte("fix"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "x", 2),
		nil,
	},
	// INVALID / missing colon after type test
	{
		"invalid-after-valid-type-test",
		[]byte("test"),
		false,
		nil,
		nil, // no partial result because it is not a minimal valid commit message
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "t", 3),
		nil,
	},
	// INVALID / invalid type (2 char) + colon
	{
		"invalid-type-2-char-colon",
		[]byte("ch:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 2),
		nil,
	},
	// INVALID / invalid type (3 char) + colon
	{
		"invalid-type-3-char-colon",
		[]byte("sty:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, ":", 3),
		nil,
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message",
		[]byte("fix: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	{
		"valid-minimal-commit-message-with-patch.bump",
		[]byte("fix: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			TypeConfig:  1,
		},
		"",
		&PatchVersion,
	},
	// VALID / minimal commit message with minor bump
	{
		"valid-minimal-commit-message-with-minor.bump",
		[]byte("feat: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "w",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "w",
			TypeConfig:  1,
		},
		"",
		&MinorVersion,
	},
	// VALID / minimal commit message with major bump
	{
		"valid-minimal-commit-message-with-major.bump",
		[]byte("fix!: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			Exclamation: true,
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "w",
			Exclamation: true,
			TypeConfig:  1,
		},
		"",
		&MajorVersion,
	},
	// VALID / minimal commit message
	{
		"valid-minimal-commit-message-style",
		[]byte("style: CSS skillz"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "style",
			Description: "CSS skillz",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "style",
			Description: "CSS skillz",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / minimal commit message with uppercase type
	{
		"valid-minimal-commit-message-style-uppercase-type",
		[]byte("STYLE: CSS skillz"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "style",
			Description: "CSS skillz",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "style",
			Description: "CSS skillz",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-3-chars",
		[]byte("fix>"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ">", 3),
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-4-chars",
		[]byte("perf?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 4),
		nil,
	},
	// INVALID / missing colon after valid commit message type
	{
		"missing-colon-after-type-5-chars",
		[]byte("build?"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "?", 5),
		nil,
	},
	// VALID / type + scope + description
	{
		"valid-with-scope",
		[]byte("refactor(xyz): ccc"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "refactor",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "refactor",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / uppercase type + scope + description
	{
		"valid-with-scope-uppercase-type",
		[]byte("REFACTOR(xyz): ccc"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "refactor",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "refactor",
			Scope:       cctesting.StringAddress("xyz"),
			Description: "ccc",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / type + scope + multiple whitespaces + description
	{
		"valid-with-scope-multiple-whitespaces",
		[]byte("fix(aaa):          bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope",
		[]byte("fix(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / type + scope + breaking + description
	{
		"valid-breaking-with-scope-feat",
		[]byte("feat(aaa)!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("aaa"),
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / empty scope is ignored
	{
		"valid-empty-scope-is-ignored",
		[]byte("fix(): bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / type + empty scope + breaking + description
	{
		"valid-breaking-with-empty-scope",
		[]byte("fix()!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / type + breaking + description
	{
		"valid-breaking-without-scope",
		[]byte("fix!: bbb"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "bbb",
			Exclamation: true,
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// INVALID / missing whitespace after colon (with breaking)
	{
		"invalid-missing-ws-after-colon-with-breaking",
		[]byte("fix!:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 5),
		nil,
	},
	// INVALID / missing whitespace after colon with scope
	{
		"invalid-missing-ws-after-colon-with-scope",
		[]byte("fix(x):a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 7),
		nil,
	},
	// INVALID / missing whitespace after colon with empty scope
	{
		"invalid-missing-ws-after-colon-with-empty-scope",
		[]byte("fix():a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 6),
		nil,
	},
	// INVALID / missing whitespace after colon
	{
		"invalid-missing-ws-after-colon",
		[]byte("fix:a"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "a", 4),
		nil,
	},
	// INVALID / invalid after valid type and scope
	{
		"invalid-after-valid-type-and-scope",
		[]byte("test(scope)"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ")", 10),
		nil,
	},
	// INVALID / invalid initial character
	{
		"invalid-initial-character",
		[]byte("(type: a description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, "(", 0),
		nil,
	},
	// INVALID / invalid second character
	{
		"invalid-second-character",
		[]byte("c description"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrType+ColumnPositionTemplate, " ", 1),
		nil,
	},
	// INVALID / invalid after valid type, scope, and breaking
	{
		"invalid-after-valid-type-scope-and-breaking",
		[]byte("test(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 11),
		nil,
	},
	// INVALID / invalid after valid mixed-case type, scope, and breaking
	{
		"invalid-after-valid-mixed-case-type-scope-and-breaking",
		[]byte("Test(scope)!"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, "!", 11),
		nil,
	},
	// INVALID / invalid after valid type, scope, and colon
	{
		"invalid-after-valid-type-scope-and-colon",
		[]byte("fix(scope):"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, and colon
	{
		"invalid-after-valid-type-scope-breaking-and-colon",
		[]byte("ci(scope)!:"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrEarly+ColumnPositionTemplate, ":", 10),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-space
	{
		"invalid-after-valid-type-scope-breaking-colon-and-space",
		[]byte("revert(scope)!: "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 16),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("ci(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 13),
		nil,
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("chore(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 6),
		nil,
	},
	// INVALID / incomplete scope
	{
		"invalid-incomplete-scope",
		[]byte("fix(scope"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScopeIncomplete+ColumnPositionTemplate, "e", 9),
		nil,
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("perf(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 6),
		nil,
	},
	// INVALID / double right parentheses in place of an exclamation, or a colon
	{
		"invalid-double-right-parentheses-scope",
		[]byte("fix(a))"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, ")", 6),
		nil,
	},
	// INVALID / new left parentheses after valid scope
	{
		"invalid-new-left-parentheses-after-valid-scope",
		[]byte("build(az)("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "(", 9),
		nil,
	},
	// INVALID / newline rather than whitespace in description
	{
		"invalid-newline-rather-than-whitespace-description",
		[]byte("perf(ax):\x0A description on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescriptionInit+ColumnPositionTemplate, "\n", 9),
		nil,
	},
	// INVALID / newline after whitespace in description
	{
		"invalid-newline-after-whitespace-description",
		[]byte("feat(az): \x0Adescription on newline"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrNewline+ColumnPositionTemplate, 11),
		nil,
	},
	// INVALID / newline in the description
	// VALID / newline in description ignored in best effort mode
	{
		"invalid-newline-in-description",
		[]byte("feat(ap): new\x0Aline"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("ap"),
			Description: "new",
			TypeConfig:  1,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 14),
		nil,
	},
	// INVALID / newline in the description
	// VALID / newline in description ignored in best effort mode
	{
		"invalid-newline-in-description-2",
		[]byte("perf(at)!: rrr\x0Al"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "perf",
			Scope:       cctesting.StringAddress("at"),
			Exclamation: true,
			Description: "rrr",
			TypeConfig:  1,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
		nil,
	},
	// INVALID / newline in the description
	// VALID / until the newline
	{
		"description-ending-with-single-newline",
		[]byte("perf(at)!: rrr\x0A"),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "perf",
			Scope:       cctesting.StringAddress("at"),
			Exclamation: true,
			Description: "rrr",
			TypeConfig:  1,
		},
		fmt.Sprintf(ErrMissingBlankLineAtBeginning+ColumnPositionTemplate, 15),
		nil,
	},
	// VALID
	{
		"valid-with-multiline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details

on typos fixed.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details\n\non typos fixed."),
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID
	{
		"valid-with-singleline-body",
		[]byte(`fix: correct minor typos in code

see the issue for details.`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct minor typos in code",
			Body:        cctesting.StringAddress("see the issue for details."),
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID
	{
		"valid-with-empty-body",
		[]byte(`fix: correct something

`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID
	{
		"valid-with-multiple-blank-lines-body",
		[]byte(`fix: correct something



`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "correct something",
			TypeConfig:  1,
		},
		"",
		nil,
	},
}

var testCasesForFreeFormTypes = []testCase{
	// VALID / minimal commit message with unknown bump
	{
		"valid-minimal-commit-message-with-unknown-bump",
		[]byte("new: w"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Description: "w",
			TypeConfig:  3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "new",
			Description: "w",
			TypeConfig:  3,
		},
		"",
		&UnknownVersion,
	},
	// VALID / multi-line body (with blank lines) and multiple signed-off-by trailers
	{
		"valid-kernel-commit-multiline-body-with-blank-lines-and-multiple-signed-off-by-trailers",
		[]byte(`kconfig: highlight xconfig 'comment' lines with '***'

Mark Kconfig "comment" lines with "*** <commentstring> ***"
so that it is clear that these lines are comments and not some
kconfig item that cannot be modified.

This is helpful in some menus to be able to provide a menu
"sub-heading" for groups of similar config items.

This also makes the comments be presented in a way that is
similar to menuconfig and nconfig.

Signed-off-by: Randy Dunlap <rdunlap@infradead.org>
Signed-off-by: Masahiro Yamada <masahiroy@kernel.org>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "kconfig",
			Description: "highlight xconfig 'comment' lines with '***'",
			Body: cctesting.StringAddress(`Mark Kconfig "comment" lines with "*** <commentstring> ***"
so that it is clear that these lines are comments and not some
kconfig item that cannot be modified.

This is helpful in some menus to be able to provide a menu
"sub-heading" for groups of similar config items.

This also makes the comments be presented in a way that is
similar to menuconfig and nconfig.`),
			Footers: map[string][]string{
				"signed-off-by": {
					"Randy Dunlap <rdunlap@infradead.org>",
					"Masahiro Yamada <masahiroy@kernel.org>",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "kconfig",
			Description: "highlight xconfig 'comment' lines with '***'",
			Body: cctesting.StringAddress(`Mark Kconfig "comment" lines with "*** <commentstring> ***"
so that it is clear that these lines are comments and not some
kconfig item that cannot be modified.

This is helpful in some menus to be able to provide a menu
"sub-heading" for groups of similar config items.

This also makes the comments be presented in a way that is
similar to menuconfig and nconfig.`),
			Footers: map[string][]string{
				"signed-off-by": {
					"Randy Dunlap <rdunlap@infradead.org>",
					"Masahiro Yamada <masahiroy@kernel.org>",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / multi-line body (with blank lines and non alphanumberic character after a blank line) and multiple different trailers
	{
		"valid-kernel-commit-long-body-with-non-trailer-start-chars-after-blank-line-and-multiple-different-trailers",
		[]byte(`bpf: fix buggy r0 retval refinement for tracing helpers

See the glory details in 100605035e15 ("bpf: Verifier, do_refine_retval_range
may clamp umin to 0 incorrectly") for why 849fa50662fb ("bpf/verifier: refine
retval R0 state for bpf_get_stack helper") is buggy. The whole series however
is not suitable for stable since it adds significant amount [0] of verifier
complexity in order to add 32bit subreg tracking. Something simpler is needed.

Unfortunately, reverting 849fa50662fb ("bpf/verifier: refine retval R0 state
for bpf_get_stack helper") or just cherry-picking 100605035e15 ("bpf: Verifier,
do_refine_retval_range may clamp umin to 0 incorrectly") is not an option since
it will break existing tracing programs badly (at least those that are using
bpf_get_stack() and bpf_probe_read_str() helpers). Not fixing it in stable is
also not an option since on 4.19 kernels an error will cause a soft-lockup due
to hitting dead-code sanitized branch since we don't hard-wire such branches
in old kernels yet. But even then for 5.x 849fa50662fb ("bpf/verifier: refine
retval R0 state for bpf_get_stack helper") would cause wrong bounds on the
verifier simluation when an error is hit.

In one of the earlier iterations of mentioned patch series for upstream there
was the concern that just using smax_value in do_refine_retval_range() would
nuke bounds by subsequent <<32 >>32 shifts before the comparison against 0 [1]
which eventually led to the 32bit subreg tracking in the first place. While I
initially went for implementing the idea [1] to pattern match the two shift
operations, it turned out to be more complex than actually needed, meaning, we
could simply treat do_refine_retval_range() similarly to how we branch off
verification for conditionals or under speculation, that is, pushing a new
reg state to the stack for later verification. This means, instead of verifying
the current path with the ret_reg in [S32MIN, msize_max_value] interval where
later bounds would get nuked, we split this into two: i) for the success case
where ret_reg can be in [0, msize_max_value], and ii) for the error case with
ret_reg known to be in interval [S32MIN, -1]. Latter will preserve the bounds
during these shift patterns and can match reg < 0 test. test_progs also succeed
with this approach.

[0] https://lore.kernel.org/bpf/158507130343.15666.8018068546764556975.stgit@john-Precision-5820-Tower/
[1] https://lore.kernel.org/bpf/158015334199.28573.4940395881683556537.stgit@john-XPS-13-9370/T/#m2e0ad1d5949131014748b6daa48a3495e7f0456d

Fixes: 849fa50662fb ("bpf/verifier: refine retval R0 state for bpf_get_stack helper")
Reported-by: Lorenzo Fontana <fontanalorenz@gmail.com>
Reported-by: Leonardo Di Donato <leodidonato@gmail.com>
Reported-by: John Fastabend <john.fastabend@gmail.com>
Signed-off-by: Daniel Borkmann <daniel@iogearbox.net>
Acked-by: Alexei Starovoitov <ast@kernel.org>
Acked-by: John Fastabend <john.fastabend@gmail.com>
Tested-by: John Fastabend <john.fastabend@gmail.com>
Tested-by: Lorenzo Fontana <fontanalorenz@gmail.com>
Tested-by: Leonardo Di Donato <leodidonato@gmail.com>
Signed-off-by: Greg Kroah-Hartman <gregkh@linuxfoundation.org>`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "bpf",
			Description: "fix buggy r0 retval refinement for tracing helpers",
			Body: cctesting.StringAddress(`See the glory details in 100605035e15 ("bpf: Verifier, do_refine_retval_range
may clamp umin to 0 incorrectly") for why 849fa50662fb ("bpf/verifier: refine
retval R0 state for bpf_get_stack helper") is buggy. The whole series however
is not suitable for stable since it adds significant amount [0] of verifier
complexity in order to add 32bit subreg tracking. Something simpler is needed.

Unfortunately, reverting 849fa50662fb ("bpf/verifier: refine retval R0 state
for bpf_get_stack helper") or just cherry-picking 100605035e15 ("bpf: Verifier,
do_refine_retval_range may clamp umin to 0 incorrectly") is not an option since
it will break existing tracing programs badly (at least those that are using
bpf_get_stack() and bpf_probe_read_str() helpers). Not fixing it in stable is
also not an option since on 4.19 kernels an error will cause a soft-lockup due
to hitting dead-code sanitized branch since we don't hard-wire such branches
in old kernels yet. But even then for 5.x 849fa50662fb ("bpf/verifier: refine
retval R0 state for bpf_get_stack helper") would cause wrong bounds on the
verifier simluation when an error is hit.

In one of the earlier iterations of mentioned patch series for upstream there
was the concern that just using smax_value in do_refine_retval_range() would
nuke bounds by subsequent <<32 >>32 shifts before the comparison against 0 [1]
which eventually led to the 32bit subreg tracking in the first place. While I
initially went for implementing the idea [1] to pattern match the two shift
operations, it turned out to be more complex than actually needed, meaning, we
could simply treat do_refine_retval_range() similarly to how we branch off
verification for conditionals or under speculation, that is, pushing a new
reg state to the stack for later verification. This means, instead of verifying
the current path with the ret_reg in [S32MIN, msize_max_value] interval where
later bounds would get nuked, we split this into two: i) for the success case
where ret_reg can be in [0, msize_max_value], and ii) for the error case with
ret_reg known to be in interval [S32MIN, -1]. Latter will preserve the bounds
during these shift patterns and can match reg < 0 test. test_progs also succeed
with this approach.

[0] https://lore.kernel.org/bpf/158507130343.15666.8018068546764556975.stgit@john-Precision-5820-Tower/
[1] https://lore.kernel.org/bpf/158015334199.28573.4940395881683556537.stgit@john-XPS-13-9370/T/#m2e0ad1d5949131014748b6daa48a3495e7f0456d`),
			Footers: map[string][]string{
				"acked-by": {
					"Alexei Starovoitov <ast@kernel.org>",
					"John Fastabend <john.fastabend@gmail.com>",
				},
				"fixes": {
					"849fa50662fb (\"bpf/verifier: refine retval R0 state for bpf_get_stack helper\")",
				},
				"reported-by": {
					"Lorenzo Fontana <fontanalorenz@gmail.com>",
					"Leonardo Di Donato <leodidonato@gmail.com>",
					"John Fastabend <john.fastabend@gmail.com>",
				},
				"signed-off-by": {
					"Daniel Borkmann <daniel@iogearbox.net>",
					"Greg Kroah-Hartman <gregkh@linuxfoundation.org>",
				},
				"tested-by": {
					"John Fastabend <john.fastabend@gmail.com>",
					"Lorenzo Fontana <fontanalorenz@gmail.com>",
					"Leonardo Di Donato <leodidonato@gmail.com>",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "bpf",
			Description: "fix buggy r0 retval refinement for tracing helpers",
			Body: cctesting.StringAddress(`See the glory details in 100605035e15 ("bpf: Verifier, do_refine_retval_range
may clamp umin to 0 incorrectly") for why 849fa50662fb ("bpf/verifier: refine
retval R0 state for bpf_get_stack helper") is buggy. The whole series however
is not suitable for stable since it adds significant amount [0] of verifier
complexity in order to add 32bit subreg tracking. Something simpler is needed.

Unfortunately, reverting 849fa50662fb ("bpf/verifier: refine retval R0 state
for bpf_get_stack helper") or just cherry-picking 100605035e15 ("bpf: Verifier,
do_refine_retval_range may clamp umin to 0 incorrectly") is not an option since
it will break existing tracing programs badly (at least those that are using
bpf_get_stack() and bpf_probe_read_str() helpers). Not fixing it in stable is
also not an option since on 4.19 kernels an error will cause a soft-lockup due
to hitting dead-code sanitized branch since we don't hard-wire such branches
in old kernels yet. But even then for 5.x 849fa50662fb ("bpf/verifier: refine
retval R0 state for bpf_get_stack helper") would cause wrong bounds on the
verifier simluation when an error is hit.

In one of the earlier iterations of mentioned patch series for upstream there
was the concern that just using smax_value in do_refine_retval_range() would
nuke bounds by subsequent <<32 >>32 shifts before the comparison against 0 [1]
which eventually led to the 32bit subreg tracking in the first place. While I
initially went for implementing the idea [1] to pattern match the two shift
operations, it turned out to be more complex than actually needed, meaning, we
could simply treat do_refine_retval_range() similarly to how we branch off
verification for conditionals or under speculation, that is, pushing a new
reg state to the stack for later verification. This means, instead of verifying
the current path with the ret_reg in [S32MIN, msize_max_value] interval where
later bounds would get nuked, we split this into two: i) for the success case
where ret_reg can be in [0, msize_max_value], and ii) for the error case with
ret_reg known to be in interval [S32MIN, -1]. Latter will preserve the bounds
during these shift patterns and can match reg < 0 test. test_progs also succeed
with this approach.

[0] https://lore.kernel.org/bpf/158507130343.15666.8018068546764556975.stgit@john-Precision-5820-Tower/
[1] https://lore.kernel.org/bpf/158015334199.28573.4940395881683556537.stgit@john-XPS-13-9370/T/#m2e0ad1d5949131014748b6daa48a3495e7f0456d`),
			Footers: map[string][]string{
				"acked-by": {
					"Alexei Starovoitov <ast@kernel.org>",
					"John Fastabend <john.fastabend@gmail.com>",
				},
				"fixes": {
					"849fa50662fb (\"bpf/verifier: refine retval R0 state for bpf_get_stack helper\")",
				},
				"reported-by": {
					"Lorenzo Fontana <fontanalorenz@gmail.com>",
					"Leonardo Di Donato <leodidonato@gmail.com>",
					"John Fastabend <john.fastabend@gmail.com>",
				},
				"signed-off-by": {
					"Daniel Borkmann <daniel@iogearbox.net>",
					"Greg Kroah-Hartman <gregkh@linuxfoundation.org>",
				},
				"tested-by": {
					"John Fastabend <john.fastabend@gmail.com>",
					"Lorenzo Fontana <fontanalorenz@gmail.com>",
					"Leonardo Di Donato <leodidonato@gmail.com>",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / type containing slash
	{
		"valid-kernel-commit-type-containing-slash",
		[]byte(`selftests/bpf: Fix core_reloc test runner

Fix failed tests checks in core_reloc test runner, which allowed failing tests
to pass quietly. Also add extra check to make sure that expected to fail test cases with
invalid names are caught as test failure anyway, as this is not an expected
failure mode. Also fix mislabeled probed vs direct bitfield test cases.

Fixes: 124a892d1c41 ("selftests/bpf: Test TYPE_EXISTS and TYPE_SIZE CO-RE relocations")
Reported-by: Lorenz Bauer <lmb@cloudflare.com>
Signed-off-by: Andrii Nakryiko <andrii@kernel.org>
Signed-off-by: Alexei Starovoitov <ast@kernel.org>
Acked-by: Lorenz Bauer <lmb@cloudflare.com>
Link: https://lore.kernel.org/bpf/20210426192949.416837-6-andrii@kernel.org`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "selftests/bpf",
			Description: "Fix core_reloc test runner",
			Body: cctesting.StringAddress(`Fix failed tests checks in core_reloc test runner, which allowed failing tests
to pass quietly. Also add extra check to make sure that expected to fail test cases with
invalid names are caught as test failure anyway, as this is not an expected
failure mode. Also fix mislabeled probed vs direct bitfield test cases.`),
			Footers: map[string][]string{
				"fixes": {
					"124a892d1c41 (\"selftests/bpf: Test TYPE_EXISTS and TYPE_SIZE CO-RE relocations\")",
				},
				"reported-by": {
					"Lorenz Bauer <lmb@cloudflare.com>",
				},
				"signed-off-by": {
					"Andrii Nakryiko <andrii@kernel.org>",
					"Alexei Starovoitov <ast@kernel.org>",
				},
				"acked-by": {
					"Lorenz Bauer <lmb@cloudflare.com>",
				},
				"link": {
					"https://lore.kernel.org/bpf/20210426192949.416837-6-andrii@kernel.org",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "selftests/bpf",
			Description: "Fix core_reloc test runner",
			Body: cctesting.StringAddress(`Fix failed tests checks in core_reloc test runner, which allowed failing tests
to pass quietly. Also add extra check to make sure that expected to fail test cases with
invalid names are caught as test failure anyway, as this is not an expected
failure mode. Also fix mislabeled probed vs direct bitfield test cases.`),
			Footers: map[string][]string{
				"fixes": {
					"124a892d1c41 (\"selftests/bpf: Test TYPE_EXISTS and TYPE_SIZE CO-RE relocations\")",
				},
				"reported-by": {
					"Lorenz Bauer <lmb@cloudflare.com>",
				},
				"signed-off-by": {
					"Andrii Nakryiko <andrii@kernel.org>",
					"Alexei Starovoitov <ast@kernel.org>",
				},
				"acked-by": {
					"Lorenz Bauer <lmb@cloudflare.com>",
				},
				"link": {
					"https://lore.kernel.org/bpf/20210426192949.416837-6-andrii@kernel.org",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / colon and space separator in the description
	{
		"valid-colon-separator-in-description",
		[]byte(`bpf: selftests: Add kfunc_call test

Signed-off-by: Martin KaFai Lau <kafai@fb.com>
Signed-off-by: Alexei Starovoitov <ast@kernel.org>
Link: https://lore.kernel.org/bpf/20210325015252.1551395-1-kafai@fb.com`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "bpf",
			Description: "selftests: Add kfunc_call test",
			Footers: map[string][]string{
				"signed-off-by": {
					"Martin KaFai Lau <kafai@fb.com>",
					"Alexei Starovoitov <ast@kernel.org>",
				},
				"link": {
					"https://lore.kernel.org/bpf/20210325015252.1551395-1-kafai@fb.com",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "bpf",
			Description: "selftests: Add kfunc_call test",
			Footers: map[string][]string{
				"signed-off-by": {
					"Martin KaFai Lau <kafai@fb.com>",
					"Alexei Starovoitov <ast@kernel.org>",
				},
				"link": {
					"https://lore.kernel.org/bpf/20210325015252.1551395-1-kafai@fb.com",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / free form type containing comma and space
	{
		"valid-type-containing-comma-and-space",
		[]byte(`bpf, selftests: test_maps generating unrecognized data section`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "bpf, selftests",
			Description: "test_maps generating unrecognized data section",
			TypeConfig:  3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "bpf, selftests",
			Description: "test_maps generating unrecognized data section",
			TypeConfig:  3,
		},
		"",
		nil,
	},
	// VALID / valid type (uppercase) + description starting with type-like string
	{
		"valid-type-uppercase-and-description-starting-with-type-like-string",
		[]byte(`KVM: nVMX: Truncate base/index GPR value on address calc in !64-bit`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "kvm",
			Description: "nVMX: Truncate base/index GPR value on address calc in !64-bit",
			TypeConfig:  3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "kvm",
			Description: "nVMX: Truncate base/index GPR value on address calc in !64-bit",
			TypeConfig:  3,
		},
		"",
		nil,
	},
	// VALID / free form type with scope
	{
		"valid-free-form-type-with-scope",
		[]byte(`KVM(nVMX): Truncate base/index GPR value on address calc in !64-bit`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "kvm",
			Scope:       cctesting.StringAddress("nvmx"),
			Description: "Truncate base/index GPR value on address calc in !64-bit",
			TypeConfig:  3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "kvm",
			Scope:       cctesting.StringAddress("nvmx"),
			Description: "Truncate base/index GPR value on address calc in !64-bit",
			TypeConfig:  3,
		},
		"",
		nil,
	},
	// INVALID / text after well-formed scope
	{
		"invalid-text-after-well-formed-scope",
		[]byte(`some(scope)text: aaaa`),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrColon+ColumnPositionTemplate, "t", 11),
		nil,
	},
	// INVALID / invalid after valid type, scope, breaking, colon, and white-spaces
	{
		"invalid-after-valid-type-scope-breaking-colon-and-spaces",
		[]byte("fix(scope)!:  "),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrDescription+ColumnPositionTemplate, " ", 14),
		nil,
	},
	// INVALID / double left parentheses in scope
	{
		"invalid-double-left-parentheses-scope",
		[]byte("fix(("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 4),
		nil,
	},
	// INVALID / incomplete scope
	{
		"invalid-incomplete-scope",
		[]byte("fix(scope"),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScopeIncomplete+ColumnPositionTemplate, "e", 9),
		nil,
	},
	// INVALID / double left parentheses in scope after valid character
	{
		"invalid-double-left-parentheses-scope-after-valid-character",
		[]byte("fix(a("),
		false,
		nil,
		nil,
		fmt.Sprintf(ErrScope+ColumnPositionTemplate, "(", 5),
		nil,
	},
	// VALID / breaking free form type with scope
	{
		"valid-breaking-free-form-type-with-scope",
		[]byte(`some(scope)!: breaking desc`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "some",
			Scope:       cctesting.StringAddress("scope"),
			Description: "breaking desc",
			Exclamation: true,
			TypeConfig:  3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "some",
			Scope:       cctesting.StringAddress("scope"),
			Description: "breaking desc",
			Exclamation: true,
			TypeConfig:  3,
		},
		"",
		nil,
	},
	// VALID / breaking change trailer
	{
		"valid-breaking-change-space-trailer",
		[]byte(`fix: description

BREAKING CHANGE: APIs`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / breaking-change trailer
	{
		"valid-breaking-change-trailer",
		[]byte(`fix: description

BREAKING-CHANGE: APIs`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / breaking change trailer before other trailers
	{
		"valid-breaking-change-space-trailer-before-others",
		[]byte(`fix: description

BREAKING CHANGE: APIs
Acked-by: Leo Di Donato`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
				"acked-by": {
					"Leo Di Donato",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
				"acked-by": {
					"Leo Di Donato",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / breaking change trailer after trailers
	{
		"valid-breaking-change-space-trailer-after-others",
		[]byte(`fix: description


Acked-by: Leo Di Donato
BREAKING CHANGE: APIs`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
				"acked-by": {
					"Leo Di Donato",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
				"acked-by": {
					"Leo Di Donato",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / breaking change trailer after blank lines and other trailers
	{
		"valid-breaking-change-space-trailer-after-blanklines-and-others",
		[]byte(`fix: description


Acked-by: Leo Di Donato


BREAKING CHANGE: APIs`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
				"acked-by": {
					"Leo Di Donato",
				},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"breaking-change": {
					"APIs",
				},
				"acked-by": {
					"Leo Di Donato",
				},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / invalid BREAKING CHANGE trailer separator after body
	// Note that because of the wrong separator (#) the BREAKING CHANGE trailer gets discarded as a footer component and captured as body content
	{
		"valid-breaking-change-invalid-separator-after-body",
		[]byte(`fix: description

Some text.

BREAKING CHANGE #5`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Body:        cctesting.StringAddress("Some text.\n\nBREAKING CHANGE #5"),
			TypeConfig:  3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Body:        cctesting.StringAddress("Some text.\n\nBREAKING CHANGE #5"),
			TypeConfig:  3,
		},
		"",
		nil,
	},
	// INVALID / invalid BREAKING CHANGE trailer separator after valid trailers
	// VALID / until the last valid footer trailer
	{
		"invalid-breaking-change-invalid-separator-after-trailers",
		[]byte(`fix: description

Tested-by: Leo
BREAKING CHANGE #5`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailer+ColumnPositionTemplate, " ", 48),
		nil,
	},
	// INVALID / incomplete BREAKING CHANGE trailer
	// VALID / until the last valid footer trailer
	{
		"invalid-breaking-change-incomplete-after-trailers",
		[]byte(`fix: description

Tested-by: Leo
BREAKING CHANG: XYZ`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailer+ColumnPositionTemplate, ":", 47),
		nil,
	},
	// INVALID / lowercase (space separated) BREAKING CHANGE trailer
	// VALID / until the last valid footer trailer
	{
		"invalid-lowercase-space-separated-breaking-change-after-trailers",
		[]byte(`fix: description

Tested-by: Leo
breaking change: xyz`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailer+ColumnPositionTemplate, "c", 42),
		nil,
	},
	// INVALID / illegal trailer after valid trailer
	// VALID / until the last valid footer trailer
	{
		"invalid-illegal-trailer-after-valid-trailer",
		[]byte(`fix: description

Tested-by: Leo
!`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailer+ColumnPositionTemplate, "!", 33),
		nil,
	},
	// INVALID / illegal trailer after valid trailer with an ending newline
	// VALID / until the last valid footer trailer
	{
		"invalid-illegal-trailer-after-valid-trailer-with-ending-newline",
		[]byte(`fix: description

Tested-by: Leo
a
`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailer+ColumnPositionTemplate, "\n", 34),
		nil,
	},
	// INVALID / incomplete trailer after valid trailer with an ending newline
	// VALID / until the last valid footer trailer
	{
		"invalid-incomplete-trailer-after-valid-trailer-with-ending-newline",
		[]byte(`fix: description

Tested-by: Leo
a`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailerIncomplete+ColumnPositionTemplate, "a", 34),
		nil,
	},
	// INVALID / incomplete trailer after valid trailer
	// VALID / until the last valid footer trailer
	{
		"invalid-incomplete-trailer-after-valid-trailer",
		[]byte(`fix: description

Tested-by: Leo
X-
Another-trailer: x`),
		false,
		nil,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo"},
			},
			TypeConfig: 3,
		},
		fmt.Sprintf(ErrTrailer+ColumnPositionTemplate, "\n", 35),
		nil,
	},
}
