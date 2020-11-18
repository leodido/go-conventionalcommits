package conventionalcommits

// BestEfforter is an interface that wraps the HasBestEffort method.
type BestEfforter interface {
	WithBestEffort()
	HasBestEffort() bool
}

// Machine represent a FSM able to parse a conventional commit and return it in an structured way.
type Machine interface {
	Parse(input []byte) (Message, error)
	BestEfforter
}

// MachineOption represents the type of option setters for Machine instances.
type MachineOption func(m Machine) Machine

// Message represent a conventional commit message.
type Message interface {
	Ok() bool
	IsBreakingChange() bool
}

// Minimal represent a base struct for Conventional Commit messages.
type Minimal struct {
	Type        string
	Description string
	Scope       *string // can be nil
	Exclamation bool
}

// Ok tells whether the receiving commit message is well-formed or not.
//
// A minimally well-formed commit message has at least a valid type and a non empty description.
func (m *Minimal) Ok() bool {
	// todo > constraint type to a set of values
	return m.Type != "" && m.Description != ""
}

// IsBreakingChange tells whether the receiving commit message represents a breaking change or not.
func (m *Minimal) IsBreakingChange() bool {
	return m.Exclamation
}
