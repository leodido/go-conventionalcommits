package slim

import (
	"github.com/leodido/go-conventionalcommits"
)


type conventionalCommit struct {
	_type       string
	descr       string
	scope       string
	exclamation bool
	body        string
}

func (c *conventionalCommit) minimal() bool {
	return c._type != "" && c.descr != ""
}

func (c *conventionalCommit) export() conventionalcommits.Message {
	out := &conventionalcommits.ConventionalCommit{}
	out.Exclamation = c.exclamation
	out.Type = c._type
	out.Description = c.descr
	if c.scope != "" {
		out.Scope = &c.scope
	}
	if c.body != "" {
		out.Body = &c.body
	}

	return out
}
