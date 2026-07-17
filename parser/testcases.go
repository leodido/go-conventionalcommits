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
	// VALID / terminal line ending
	{
		"description-ending-with-single-newline",
		[]byte("feat(az)!: bla\x0A"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  0,
		},
		"",
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
	// --- Reproducers for issue #38 (trailer-looking lines abort parsing).
	// See https://github.com/leodido/go-conventionalcommits/issues/38

	// VALID / fake trailer at the start of body, followed by body prose.
	{
		"valid-fake-trailer-then-body-paragraph-issue38",
		[]byte("feat: x\n\nFixes #15\n\nLorem ipsum dolor sit amet"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			TypeConfig:  0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			TypeConfig:  0,
		},
		"",
		nil,
	},
	// VALID / verbatim issue #38 input.
	{
		"valid-issue-38-verbatim",
		[]byte("feat: some thing (hz/fl!144)\n\nFixes #15\n\nLorem ipsum dolor sit amet\n\nBREAKING CHANGE: Some explanation\nReviewed-by: XX <xx@example.com>"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "some thing (hz/fl!144)",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			Footers: map[string][]string{
				"breaking-change": {"Some explanation"},
				"reviewed-by":     {"XX <xx@example.com>"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "some thing (hz/fl!144)",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			Footers: map[string][]string{
				"breaking-change": {"Some explanation"},
				"reviewed-by":     {"XX <xx@example.com>"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / body paragraph whose LAST line is trailer-shaped, followed
	// by a real trailer block. Must not regress.
	{
		"valid-body-ending-trailer-shape-then-real-trailer",
		[]byte("feat: x\n\nbody1\nFake: trail\n\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nFake: trail"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nFake: trail"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / two body paragraphs each ending in a trailer-shaped line,
	// followed by a real trailer block. Must not regress.
	{
		"valid-multi-paragraph-body-each-ending-trailer-shape-then-real-trailer",
		[]byte("feat: x\n\npara1\nFake: a\n\npara2\nFake: b\n\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("para1\nFake: a\n\npara2\nFake: b"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("para1\nFake: a\n\npara2\nFake: b"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / minimised: body's last line is itself trailer-shaped
	// (`Refs: 123`) and is followed by a blank line and a real trailer
	// (`Reviewed-by: X`). The pre-scan must commit the trailer block
	// at the gap above `Reviewed-by:` rather than extending the
	// candidate run upward across the `Refs: 123` continuation. Catches
	// the second-pass review's §2 shape (body1\nRefs: 123 / Reviewed-by:)
	// which lives between the two cases above and would otherwise drift.
	{
		"valid-body-ending-trailer-shape-refs-then-real-trailer",
		[]byte("fix: x\n\nbody1\nRefs: 123\n\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nRefs: 123"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nRefs: 123"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},

	// --- Reproducers for issue #48 (multi-line trailer values per spec
	// clause 10: "A footer's value MAY contain spaces and newlines, and
	// parsing MUST terminate when the next valid footer token/separator
	// pair is observed").
	// See https://github.com/leodido/go-conventionalcommits/issues/48

	// VALID / single trailer with a multi-line value terminated by EOF.
	{
		"valid-issue-48-breaking-change-multiline-value-eof",
		[]byte("fix: x\n\nBREAKING CHANGE: this wraps\nacross two lines"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"this wraps\nacross two lines"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"this wraps\nacross two lines"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / multi-line trailer value followed by another single-line trailer.
	{
		"valid-issue-48-multiline-trailer-then-single-trailer",
		[]byte("feat: x\n\nBREAKING CHANGE: long\nexplanation continues here\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"long\nexplanation continues here"},
				"reviewed-by":     {"X"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"long\nexplanation continues here"},
				"reviewed-by":     {"X"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / multi-line trailer value where the continuation line begins
	// with whitespace. Continuation must be retained verbatim, including
	// the leading space, since a leading space disqualifies the line from
	// being a trailer_init.
	{
		"valid-issue-48-multiline-trailer-value-indented-continuation",
		[]byte("fix: x\n\nRefs: #123\n  also fixes #124"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"refs": {"#123\n  also fixes #124"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"refs": {"#123\n  also fixes #124"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / continuation line that visually resembles a trailer but
	// has a leading space (so it is NOT trailer-init-shaped) is part of
	// the value of the trailer above; the real trailer below terminates.
	{
		"valid-issue-48-fake-trailer-continuation-then-real-trailer",
		[]byte("fix: x\n\nCo-authored-by: Alice\n Co-authored-by: Bob\nReviewed-by: Carol"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"co-authored-by": {"Alice\n Co-authored-by: Bob"},
				"reviewed-by":    {"Carol"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"co-authored-by": {"Alice\n Co-authored-by: Bob"},
				"reviewed-by":    {"Carol"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / clause-8 shape: trailers separated by a blank line. The
	// trailer_val_continues guard MUST treat the current \n as the
	// first half of the blank-line gap and stop the value before it,
	// rather than consuming this \n into the value and only stopping
	// at the next one. Pins the §B.2 regression from PR #49's
	// third-pass review.
	{
		"valid-issue-48-trailers-separated-by-blank-line",
		[]byte("feat: x\n\nA: 1\n\nB: 2"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / three trailers each separated by blank lines (clause 8).
	// The leak surfaces on every trailer that is itself the last of a
	// paragraph followed by a blank line — i.e. all but the last
	// trailer in this shape.
	{
		"valid-issue-48-three-trailers-each-separated-by-blank-line",
		[]byte("feat: x\n\nA: 1\n\nB: 2\n\nC: 3"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / value followed by EXACTLY one trailing blank line at EOF.
	// Mirrors `A: 1\n\n` — the blank-line gap before EOF must NOT be
	// folded into the value. Pins the EOF-asymmetry sub-case of §B.2.
	{
		"valid-issue-48-trailer-then-trailing-blank-line-at-eof",
		[]byte("feat: x\n\nA: 1\n\n"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
			},
			TypeConfig: 0,
		},
		"",
		nil,
	},
	// VALID / Git-canonical shape: real-world commit message with a
	// `Fixes #N` paragraph followed by a `Signed-off-by:` paragraph.
	// This is the exact shape spec clause 8 documents and the one
	// PR #49's third-pass review identified as silently-broken in the
	// wild. Worth pinning verbatim.
	{
		"valid-issue-48-git-canonical-fixes-then-signed-off-by",
		[]byte("feat: add new widget\n\nFixes #123\n\nSigned-off-by: Alice <alice@example.com>"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "add new widget",
			Footers: map[string][]string{
				"fixes":         {"123"},
				"signed-off-by": {"Alice <alice@example.com>"},
			},
			TypeConfig: 0,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "add new widget",
			Footers: map[string][]string{
				"fixes":         {"123"},
				"signed-off-by": {"Alice <alice@example.com>"},
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
	// VALID / terminal line ending
	{
		"description-ending-with-single-newline",
		[]byte("docs(az)!: bla\x0A"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "docs",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  2,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "docs",
			Scope:       cctesting.StringAddress("az"),
			Exclamation: true,
			Description: "bla",
			TypeConfig:  2,
		},
		"",
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
	// VALID / terminal line ending
	{
		"description-ending-with-single-newline",
		[]byte("perf(at)!: rrr\x0A"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "perf",
			Scope:       cctesting.StringAddress("at"),
			Exclamation: true,
			Description: "rrr",
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "perf",
			Scope:       cctesting.StringAddress("at"),
			Exclamation: true,
			Description: "rrr",
			TypeConfig:  1,
		},
		"",
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
	// --- Reproducers for issue #38 (trailer-looking lines abort parsing).
	// See https://github.com/leodido/go-conventionalcommits/issues/38

	// VALID / fake trailer at the start of body, followed by body prose.
	{
		"valid-fake-trailer-then-body-paragraph-issue38",
		[]byte("feat: x\n\nFixes #15\n\nLorem ipsum dolor sit amet"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			TypeConfig:  1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			TypeConfig:  1,
		},
		"",
		nil,
	},
	// VALID / verbatim issue #38 input.
	{
		"valid-issue-38-verbatim",
		[]byte("feat: some thing (hz/fl!144)\n\nFixes #15\n\nLorem ipsum dolor sit amet\n\nBREAKING CHANGE: Some explanation\nReviewed-by: XX <xx@example.com>"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "some thing (hz/fl!144)",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			Footers: map[string][]string{
				"breaking-change": {"Some explanation"},
				"reviewed-by":     {"XX <xx@example.com>"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "some thing (hz/fl!144)",
			Body:        cctesting.StringAddress("Fixes #15\n\nLorem ipsum dolor sit amet"),
			Footers: map[string][]string{
				"breaking-change": {"Some explanation"},
				"reviewed-by":     {"XX <xx@example.com>"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	// VALID / body paragraph whose LAST line is trailer-shaped, followed
	// by a real trailer block. Must not regress.
	{
		"valid-body-ending-trailer-shape-then-real-trailer",
		[]byte("feat: x\n\nbody1\nFake: trail\n\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nFake: trail"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nFake: trail"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	// VALID / two body paragraphs each ending in a trailer-shaped line,
	// followed by a real trailer block. Must not regress.
	{
		"valid-multi-paragraph-body-each-ending-trailer-shape-then-real-trailer",
		[]byte("feat: x\n\npara1\nFake: a\n\npara2\nFake: b\n\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("para1\nFake: a\n\npara2\nFake: b"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Body:        cctesting.StringAddress("para1\nFake: a\n\npara2\nFake: b"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	// VALID / minimised: body's last line is itself trailer-shaped
	// (`Refs: 123`) and is followed by a blank line and a real trailer
	// (`Reviewed-by: X`). The pre-scan must commit the trailer block
	// at the gap above `Reviewed-by:` rather than extending the
	// candidate run upward across the `Refs: 123` continuation. Catches
	// the second-pass review's §2 shape (body1\nRefs: 123 / Reviewed-by:)
	// which lives between the two cases above and would otherwise drift.
	{
		"valid-body-ending-trailer-shape-refs-then-real-trailer",
		[]byte("fix: x\n\nbody1\nRefs: 123\n\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nRefs: 123"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Body:        cctesting.StringAddress("body1\nRefs: 123"),
			Footers: map[string][]string{
				"reviewed-by": {"X"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},

	// --- Reproducers for issue #48 (multi-line trailer values per spec
	// clause 10). See https://github.com/leodido/go-conventionalcommits/issues/48
	{
		"valid-issue-48-breaking-change-multiline-value-eof",
		[]byte("fix: x\n\nBREAKING CHANGE: this wraps\nacross two lines"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"this wraps\nacross two lines"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"this wraps\nacross two lines"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	{
		"valid-issue-48-multiline-trailer-then-single-trailer",
		[]byte("feat: x\n\nBREAKING CHANGE: long\nexplanation continues here\nReviewed-by: X"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"long\nexplanation continues here"},
				"reviewed-by":     {"X"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"breaking-change": {"long\nexplanation continues here"},
				"reviewed-by":     {"X"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	{
		"valid-issue-48-multiline-trailer-value-indented-continuation",
		[]byte("fix: x\n\nRefs: #123\n  also fixes #124"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"refs": {"#123\n  also fixes #124"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"refs": {"#123\n  also fixes #124"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	{
		"valid-issue-48-fake-trailer-continuation-then-real-trailer",
		[]byte("fix: x\n\nCo-authored-by: Alice\n Co-authored-by: Bob\nReviewed-by: Carol"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"co-authored-by": {"Alice\n Co-authored-by: Bob"},
				"reviewed-by":    {"Carol"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "x",
			Footers: map[string][]string{
				"co-authored-by": {"Alice\n Co-authored-by: Bob"},
				"reviewed-by":    {"Carol"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	// VALID / clause-8 shape: trailers separated by a blank line. The
	// trailer_val_continues guard MUST treat the current \n as the
	// first half of the blank-line gap and stop the value before it,
	// rather than consuming this \n into the value and only stopping
	// at the next one. Pins the §B.2 regression from PR #49's
	// third-pass review.
	{
		"valid-issue-48-trailers-separated-by-blank-line",
		[]byte("feat: x\n\nA: 1\n\nB: 2"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	{
		"valid-issue-48-three-trailers-each-separated-by-blank-line",
		[]byte("feat: x\n\nA: 1\n\nB: 2\n\nC: 3"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
				"b": {"2"},
				"c": {"3"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	{
		"valid-issue-48-trailer-then-trailing-blank-line-at-eof",
		[]byte("feat: x\n\nA: 1\n\n"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "x",
			Footers: map[string][]string{
				"a": {"1"},
			},
			TypeConfig: 1,
		},
		"",
		nil,
	},
	{
		"valid-issue-48-git-canonical-fixes-then-signed-off-by",
		[]byte("feat: add new widget\n\nFixes #123\n\nSigned-off-by: Alice <alice@example.com>"),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "add new widget",
			Footers: map[string][]string{
				"fixes":         {"123"},
				"signed-off-by": {"Alice <alice@example.com>"},
			},
			TypeConfig: 1,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "feat",
			Description: "add new widget",
			Footers: map[string][]string{
				"fixes":         {"123"},
				"signed-off-by": {"Alice <alice@example.com>"},
			},
			TypeConfig: 1,
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
	// VALID / per spec clause 10 a multi-line trailer value MAY contain
	// arbitrary text including lines that look like a malformed
	// BREAKING-CHANGE trailer; parsing only terminates at a real
	// trailer line, a blank line, or EOF. See issue #48.
	{
		"valid-issue-48-trailer-then-non-trailer-breaking-change-line",
		[]byte(`fix: description

Tested-by: Leo
BREAKING CHANGE #5`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\nBREAKING CHANGE #5"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\nBREAKING CHANGE #5"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / per spec clause 10: a `BREAKING CHANG: XYZ` line is
	// not a trailer (typo in the token); it continues the previous
	// trailer's value. See issue #48.
	{
		"valid-issue-48-trailer-then-mistyped-breaking-change-line",
		[]byte(`fix: description

Tested-by: Leo
BREAKING CHANG: XYZ`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\nBREAKING CHANG: XYZ"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\nBREAKING CHANG: XYZ"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / per spec clause 10: lowercase `breaking change` is not
	// a real trailer token; the line continues the previous value.
	// See issue #48.
	{
		"valid-issue-48-trailer-then-lowercase-breaking-change-line",
		[]byte(`fix: description

Tested-by: Leo
breaking change: xyz`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\nbreaking change: xyz"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\nbreaking change: xyz"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / per spec clause 10: a continuation line beginning with
	// `!` is not a trailer; it is part of the previous trailer's
	// value. See issue #48.
	{
		"valid-issue-48-trailer-then-bang-only-continuation",
		[]byte(`fix: description

Tested-by: Leo
!`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\n!"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\n!"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / per spec clause 10: a continuation line of a single
	// alphabetic char (here `a`) is not trailer-shaped; the trailing
	// newline is dropped per the existing trailer-block contract.
	// See issue #48.
	{
		"valid-issue-48-trailer-then-alpha-continuation-with-trailing-newline",
		[]byte(`fix: description

Tested-by: Leo
a
`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\na"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\na"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / per spec clause 10: same as above but without the
	// trailing newline (EOF mid-continuation). See issue #48.
	{
		"valid-issue-48-trailer-then-alpha-continuation-eof",
		[]byte(`fix: description

Tested-by: Leo
a`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\na"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by": {"Leo\na"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
	// VALID / per spec clause 10: `X-` is not trailer-shaped (a
	// trailer token cannot end in a dash; trailer_init requires
	// `alnum+ (- alnum+)*`), so it folds into the previous trailer's
	// value; `Another-trailer: x` IS a real trailer line and starts
	// a new trailer. See issue #48.
	{
		"valid-issue-48-trailer-then-dash-trailing-continuation-then-real-trailer",
		[]byte(`fix: description

Tested-by: Leo
X-
Another-trailer: x`),
		true,
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by":       {"Leo\nX-"},
				"another-trailer": {"x"},
			},
			TypeConfig: 3,
		},
		&conventionalcommits.ConventionalCommit{
			Type:        "fix",
			Description: "description",
			Footers: map[string][]string{
				"tested-by":       {"Leo\nX-"},
				"another-trailer": {"x"},
			},
			TypeConfig: 3,
		},
		"",
		nil,
	},
}
