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

	// Byte offsets in the ORIGINAL input where each captured slice
	// begins. -1 when the corresponding field was never captured.
	// Used by the WithStrictUTF8 / WithStrictUTF8Body validators to
	// report the column of the first invalid byte against the
	// original input rather than against the post-parse cursor
	// (m.p == m.pe at validation time). The exported view goes
	// through strings.ToLower for type/scope, which collapses
	// invalid bytes to U+FFFD; validation reads the raw fields here
	// instead.
	typeOffset  int
	descrOffset int
	scopeOffset int
	bodyOffset  int
	// footerValueOffsets parallels footers: footerValueOffsets[k][i]
	// is the byte offset in the original input of the i-th value
	// captured for footer key k.
	footerValueOffsets map[string][]int
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
