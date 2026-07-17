// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>

// Command ccvalidate validates a single Conventional Commits header
// (the first line of a commit message, or a PR title) using this
// repository's own parser. It is the authoritative on-CI gate that
// complements the bash-based labeler in .github/workflows/label.yml.
//
// Usage:
//
//	ccvalidate "feat(parser): add foo"            # argv
//	echo "feat: add foo" | ccvalidate             # stdin (single line)
//	CC_INPUT_FILE=/path/to/title ccvalidate       # file (avoids shell quoting in CI)
//	ccvalidate --describe "feat(API)!: add foo"   # emit canonical fields
//
// Exit codes:
//
//	0  valid Conventional Commits header
//	1  invalid (parse error printed to stderr)
//	2  usage error
//
// In --describe mode, on a valid header ccvalidate writes
// canonicalised key=value lines to stdout:
//
//	type=feat
//	scope=api
//	breaking=true
//
// "type" is always lowercased — even if the parser accepted "FEAT" via
// its Ragel case-insensitive type token — so downstream consumers
// (labelers, agreement tests) get a single canonical form. "scope" is
// emitted only when present; "breaking" is "true" or "false". This
// makes ccvalidate the single source of truth: the labeler can be
// regenerated from --describe, and an agreement test can assert that
// what ccvalidate accepts is exactly what the labeler labels.
//
// The parser is run in strict mode. Best-effort mode is intentionally
// disabled: a CI gate must be authoritative.
package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	cc "github.com/leodido/go-conventionalcommits"
	"github.com/leodido/go-conventionalcommits/parser"
)

const usage = `usage: ccvalidate [--describe] <header>
       echo "<header>" | ccvalidate [--describe]
       CC_INPUT_FILE=/path/to/file ccvalidate [--describe]

Reads a Conventional Commits header from argv, stdin, or the file at
$CC_INPUT_FILE, and validates it with the conventional type set.

With --describe, on a valid header writes canonical key=value lines
(type, scope, breaking) to stdout. type is always lowercased.
`

func main() {
	if err := run(); err != nil {
		if errors.Is(err, errUsage) {
			fmt.Fprint(os.Stderr, usage)
			os.Exit(2)
		}
		fmt.Fprintf(os.Stderr, "ccvalidate: %v\n", err)
		os.Exit(1)
	}
}

var errUsage = errors.New("usage")

func run() error {
	args, describe := parseFlags(os.Args[1:])

	header, err := readHeader(args)
	if err != nil {
		return err
	}

	return validate(header, describe, os.Stdout)
}

// parseFlags strips a leading --describe flag from args. We don't use
// the stdlib `flag` package because we want predictable behaviour when
// the header itself begins with "-" (a legal Conventional Commits
// description e.g. "feat: -prefixed flag handling").
func parseFlags(args []string) ([]string, bool) {
	if len(args) > 0 && args[0] == "--describe" {
		return args[1:], true
	}

	return args, false
}

// validate parses header with the conventional type set, returning a
// non-nil error if the header is empty or fails to parse. Trailing CR
// and LF are trimmed first so file/stdin input with a trailing newline
// is accepted.
//
// When describe is true and parsing succeeds, canonical key=value
// lines are written to w:
//   - type     (always lowercased)
//   - scope    (only when the header has a scope)
//   - breaking ("true" / "false")
func validate(header string, describe bool, w io.Writer) error {
	header = strings.TrimRight(header, "\r\n")
	if header == "" {
		return fmt.Errorf("empty header")
	}

	m := parser.NewMachine(parser.WithTypes(cc.TypesConventional))
	m.WithStrictUTF8()
	msg, perr := m.Parse([]byte(header))
	if perr != nil {
		return fmt.Errorf("invalid Conventional Commits header %q: %w", header, perr)
	}

	if !describe {
		return nil
	}

	c, ok := msg.(*cc.ConventionalCommit)
	if !ok {
		// Defensive: parser.NewMachine with TypesConventional always
		// returns *ConventionalCommit on success today, but the
		// Message interface allows for other shapes. Treat any future
		// divergence as a bug rather than silently dropping fields.
		return fmt.Errorf("internal: unexpected message type %T", msg)
	}

	// Lowercase the type to canonicalise. The parser uses Ragel's 'i
	// flag on the type token (parser/machine.go.rl), which means
	// "FEAT: x" parses successfully with Type == "feat" already on
	// current code — but we lowercase explicitly so the contract holds
	// even if the Ragel grammar is regenerated without 'i.
	fmt.Fprintf(w, "type=%s\n", strings.ToLower(c.Type))
	if c.Scope != nil {
		fmt.Fprintf(w, "scope=%s\n", *c.Scope)
	}
	fmt.Fprintf(w, "breaking=%t\n", c.IsBreakingChange())

	return nil
}

func readHeader(args []string) (string, error) {
	// 1. argv
	if len(args) > 1 {
		return "", errUsage
	}
	if len(args) == 1 {
		return args[0], nil
	}

	// 2. file (used by CI to avoid shell quoting issues with
	//    attacker-controlled PR titles).
	if path := os.Getenv("CC_INPUT_FILE"); path != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("read CC_INPUT_FILE: %w", err)
		}

		return string(b), nil
	}

	// 3. stdin
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", fmt.Errorf("stat stdin: %w", err)
	}
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		// stdin is a TTY, no piped input — treat as usage error.
		return "", errUsage
	}
	b, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("read stdin: %w", err)
	}

	return string(b), nil
}
