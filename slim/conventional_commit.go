package slim

import (
	"github.com/leodido/go-conventionalcommits"
)

// ConventionalCommit represent a commit message (with no body part) that respects ConventionalCommits v1 specification.
type ConventionalCommit struct {
	conventionalcommits.Minimal
}

// Ok tells whether the receiving commit message is well-formed or not.
//
// A minimally well-formed commit message has at least a valid type and a non empty description.
func (c *ConventionalCommit) Ok() bool {
	return c.Minimal.Ok()
}

// IsBreakingChange tells whether the receiving commit message has an exclamation point to communicate it refers to a breaking change.
func (c *ConventionalCommit) IsBreakingChange() bool {
	return c.Minimal.IsBreakingChange()
}

type conventionalCommit struct {
	_type       string
	descr       string
	scope       string
	exclamation bool
}

func (c *conventionalCommit) minimal() bool {
	return c._type != "" && c.descr != ""
}

func (c *conventionalCommit) export() conventionalcommits.Message {
	out := &ConventionalCommit{}
	out.Exclamation = c.exclamation
	out.Type = c._type
	out.Description = c.descr

	if c.scope != "" {
		out.Scope = &c.scope
	}

	return out
}
