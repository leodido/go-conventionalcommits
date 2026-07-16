// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2020- Leonardo Di Donato <leodidonato@gmail.com>
package parser

import (
	"strings"

	"github.com/leodido/go-conventionalcommits"
)

type conventionalCommit struct {
	_type       string
	descr       string
	scope       string
	exclamation bool
	body        string
	footers     map[string][]string
	typeconfig  conventionalcommits.TypeConfig
}

func (c *conventionalCommit) minimal() bool {
	return c._type != "" && c.descr != ""
}

func (c *conventionalCommit) export() conventionalcommits.Message {
	out := &conventionalcommits.ConventionalCommit{}
	out.Exclamation = c.exclamation
	// Type and Scope intentionally retain the parser's historical
	// Unicode-aware lowercase normalization. Consequently, grammar
	// acceptance of high bytes does not imply byte-for-byte export:
	// cased Unicode is lowercased and ill-formed bytes become U+FFFD.
	// Description, Body, and Footers do not pass through this step.
	out.Type = strings.ToLower(c._type)
	out.Description = c.descr
	out.TypeConfig = c.typeconfig
	if c.scope != "" {
		c.scope = strings.ToLower(c.scope)
		out.Scope = &c.scope
	}
	if c.body != "" {
		// Trim suffix blank line
		if len(c.body) >= 2 && c.body[len(c.body)-2:] == "\n\n" {
			c.body = c.body[:len(c.body)-2]
		}
		out.Body = &c.body
	}
	if len(c.footers) > 0 {
		out.Footers = c.footers
	}

	return out
}
