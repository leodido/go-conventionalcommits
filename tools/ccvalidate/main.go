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
//	ccvalidate "feat(parser): add foo"           # argv
//	echo "feat: add foo" | ccvalidate            # stdin (single line)
//	CC_INPUT_FILE=/path/to/title ccvalidate      # file (avoids shell quoting in CI)
//
// Exit codes:
//
//	0  valid Conventional Commits header
//	1  invalid (parse error printed to stderr)
//	2  usage error
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

const usage = `usage: ccvalidate <header>
       echo "<header>" | ccvalidate
       CC_INPUT_FILE=/path/to/file ccvalidate

Reads a Conventional Commits header from argv, stdin, or the file at
$CC_INPUT_FILE, and validates it with the conventional type set.
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
	header, err := readHeader()
	if err != nil {
		return err
	}
	header = strings.TrimRight(header, "\r\n")
	if header == "" {
		return fmt.Errorf("empty header")
	}

	m := parser.NewMachine(
		parser.WithTypes(cc.TypesConventional),
	)
	if _, perr := m.Parse([]byte(header)); perr != nil {
		return fmt.Errorf("invalid Conventional Commits header %q: %w", header, perr)
	}

	return nil
}

func readHeader() (string, error) {
	// 1. argv
	if len(os.Args) > 2 {
		return "", errUsage
	}
	if len(os.Args) == 2 {
		return os.Args[1], nil
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
