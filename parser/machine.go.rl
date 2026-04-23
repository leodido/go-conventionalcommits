package parser

import (
	"fmt"
	"bytes"

	"github.com/leodido/go-conventionalcommits"
	"github.com/sirupsen/logrus"
)

// ColumnPositionTemplate is the template used to communicate the column where errors occur.
var ColumnPositionTemplate = ": col=%02d"

const (
	// ErrType represents an error in the type part of the commit message.
	ErrType = "illegal '%s' character in commit message type"
	// ErrTypeIncomplete represents an error when the type part of the commit message is not complete.
	ErrTypeIncomplete = "incomplete commit message type after '%s' character"
	// ErrColon is the error message that communicate that the mandatory colon after the type part of the commit message is missing.
	ErrColon = "expecting colon (':') character, got '%s' character"
	// ErrScope represents an error about illegal characters into the the scope part of the commit message.
	ErrScope = "illegal '%s' character in scope"
	// ErrScopeIncomplete represents a specific early-exit error.
	ErrScopeIncomplete = "expecting closing parentheses (')') character, got early exit after '%s' character"
	// ErrEmpty represents an error when the input is empty.
	ErrEmpty = "empty input"
	// ErrEarly represents an error when the input makes the machine exit too early.
	ErrEarly = "early exit after '%s' character"
	// ErrDescriptionInit tells the user that before of the description part a whitespace is mandatory.
	ErrDescriptionInit = "expecting at least one white-space (' ') character, got '%s' character"
	// ErrDescription tells the user that after the whitespace is mandatory a description.
	ErrDescription = "expecting a description text (without newlines) after '%s' character"
	// ErrNewline communicates an illegal newline to the user.
	ErrNewline = "illegal newline"
	// ErrMissingBlankLineAtBeginning tells the user that the a blank line is missing after the description or after the body.
	ErrMissingBlankLineAtBeginning = "missing a blank line"
	// ErrTrailer represents an error due to an unexepected character while parsing a footer trailer.
	ErrTrailer = "illegal '%s' character in trailer"
	// ErrTrailerIncomplete represent an error when a trailer is not complete.
	ErrTrailerIncomplete = "incomplete footer trailer after '%s' character"
)

%%{
machine conventionalcommits;

include common "common.rl";

# unsigned alphabet
alphtype uint8;

action mark {
	m.pb = m.p
}

# Error management

action err_empty {
	m.err = m.emitErrorWithoutCharacter(ErrEmpty)
}

action err_type {
	if m.pe > 0 {
		if m.p != m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrType)
		} else {
			// assert(m.p == m.pe)
			m.err = m.emitErrorOnPreviousCharacter(ErrTypeIncomplete)
		}
	}
}

action err_malformed_scope {
	if m.p < m.pe {
		m.err = m.emitErrorOnCurrentCharacter(ErrScope)
	}
}

action err_malformed_scope_closing {
	// assert(m.p == m.pe)
	m.err = m.emitErrorOnPreviousCharacter(ErrScopeIncomplete)
}

action err_colon {
	if m.err == nil {
		m.err = m.emitErrorOnCurrentCharacter(ErrColon)
	}
}

action err_description_init {
	if m.err == nil {
		m.err = m.emitErrorOnCurrentCharacter(ErrDescriptionInit)
	}
}

action err_description {
	if m.p < m.pe && m.data[m.p] == 10 {
		m.err = m.emitError(ErrNewline, m.p + 1)
	} else {
		// assert(m.p == m.pe)
		m.err = m.emitErrorOnPreviousCharacter(ErrDescription)
	}
}

action check_early_exit {
	if (m.p + 1) == m.pe {
		m.err = m.emitErrorOnCurrentCharacter(ErrEarly)
	}
}

action err_begin_blank_line {
	m.err = m.emitErrorWithoutCharacter(ErrMissingBlankLineAtBeginning)
}

# Setters

action set_type {
	output._type = string(m.text())
	m.emitInfo("valid commit message type", "type", output._type)
}

action set_scope {
	output.scope = string(m.text())
	m.emitInfo("valid commit message scope", "scope", output.scope)
}

action set_description {
	output.descr = string(m.text())
	m.emitInfo("valid commit message description", "description", output.descr)
}

action set_exclamation {
	output.exclamation = true
	m.emitInfo("commit message communicates a breaking change")
}

action set_body_blank_line {
	m.emitDebug("found a blank line", "pos", m.p)
}

action set_current_footer_key {
	// todo > alnum[[- ]alnum] string to lower can be more performant?
	m.currentFooterKey = string(bytes.ToLower(m.text()))
	if m.currentFooterKey == "breaking change" {
		m.currentFooterKey = "breaking-change"
	}
	m.emitDebug("possibly valid footer token", "token", m.currentFooterKey, "pos", m.p)
}

action set_footer {
	output.footers[m.currentFooterKey] = append(output.footers[m.currentFooterKey], string(m.text()))
	m.emitInfo("valid commit message footer trailer", m.currentFooterKey, string(m.text()))
}

action count_nl {
	// Increment number of newlines to use in case we're still in the body
	m.countNewlines++
	m.lastNewline = m.p
	m.emitDebug("found a newline", "pos", m.p)
}

action append_body {
	// Append newlines
	for ; m.countNewlines > 0; {
		output.body += "\n"
		m.countNewlines--
		m.emitInfo("valid commit message body content", "body", "\n")
	}
	// Append body content
	output.body += string(m.text())
	m.emitInfo("valid commit message body content", "body", string(m.text()))
}

action append_body_before_blank_line {
	// Append content to body
	m.pb++
	m.p++
	output.body += string(m.text())
	m.emitInfo("valid commit message body content", "body", string(m.text()))
	// Do not advance over the current char
	fhold;
}

# Jumps

action start_trailer_parsing {
	// shouldRedirectToBody centralizes the decision of whether the next
	// stretch of input should be parsed as body content rather than as
	// the start of a footer trailer. Keeping the predicate in Go (and
	// the fgoto in Ragel) means the action body Ragel emits at every
	// callsite stays one line of conditional + fgoto, so the generated
	// machine.go does not multiply the predicate logic across callsites.
	if m.shouldRedirectToBody() {
		// Snap the marker forward so body's append_body window starts
		// fresh on the upcoming line instead of replaying any stale
		// >mark from earlier productions (description, last footer val,
		// etc.). m.p currently points at the newline that terminated
		// the previous production, so m.p + 1 is the first byte of the
		// line we are about to treat as body content.
		m.pb = m.p + 1
		m.emitDebug("redirecting to body parsing", "pos", m.p)
		fgoto body;
	}
	m.emitDebug("try to parse a footer trailer token", "pos", m.p)
	fgoto trailer_beg;
}

action complete_trailer_parsing {
	m.emitDebug("try to parse a footer trailer value", "pos", m.p)
	fgoto trailer_end;
}

action rewind {
	if len(output.footers) == 0 {
		// Backtrack to the last marker
		// Ie., the text possibly a trailer token that is instead part of the body content
		if m.countNewlines > 0 {
			// In case new lines met while rewinding
			// advance the last marker by the number of the newlines so that they don't get parsed again
			// (they be added in the result by the body content appender)
			m.pb = m.lastNewline + 1
		}
		fexec m.pb;
		m.emitDebug("try to parse body content", "pos", m.p)
		fgoto body;
	} else {
		// A rewind happens when an error while parsing a footer trailer is encountered
		// If this is not the first footer trailer the parser can't go back to parse body content again
		// Thus, emit an error
		if m.p != m.pe {
			m.err = m.emitErrorOnCurrentCharacter(ErrTrailer)
		} else {
			// assert(m.p == m.pe)
			m.err = m.emitErrorOnPreviousCharacter(ErrTrailerIncomplete)
		}
	}
}

action blank_line_ahead { m.p + 2 < m.pe && m.data[m.p + 1] == 10 && m.data[m.p + 2] == 10 }

# trailer_val_continues reports whether the FSM, currently poised to
# consume the newline at m.p, should continue the in-progress trailer
# value past that newline. The newline is consumed iff:
#   - the line right after it is NOT a real trailer line (per the
#     pre-scan), AND
#   - the next byte pair is not a blank-line gap (which terminates the
#     trailer block per the standard FSM contract).
# This is the clause-10 termination predicate. See issue #48.
action trailer_val_continues { m.trailerValueContinues() }

# Machine definitions

minimal_types = ('fix'i | 'feat'i);

conventional_types = ('build'i | 'ci'i | 'chore'i | 'docs'i | 'feat'i | 'fix'i | 'perf'i | 'refactor'i | 'revert'i | 'style'i | 'test'i);

falco_types = ('build'i | 'ci'i | 'chore'i | 'docs'i | 'feat'i | 'fix'i | 'perf'i | 'new'i | 'revert'i | 'update'i | 'test'i | 'rule'i);

free_form_types = print+;

scope = lpar ((print* -- lpar) -- rpar) >mark %err(err_malformed_scope) %eof(err_malformed_scope_closing) %set_scope rpar;

breaking = exclamation >set_exclamation;

## todo > strict option to enforce a single whitespace?
description = ws+ >err(err_description_init) <: (any - nl)+ >mark >err(err_description) %set_description;

blank_line = nl nl >err(err_begin_blank_line) >set_body_blank_line;

trailer_tok_breaking = 'BREAKING CHANGE';

trailer_sep_breaking = colon ws+;

trailer_tok = alnum+ (dash alnum+)*;

trailer_sep = trailer_sep_breaking | (ws '#');

# A trailer value MAY contain spaces and newlines; parsing terminates
# when the next line is a real trailer (per the pre-scan), or when a
# blank line is encountered, or at EOF. The `when` guard on the nl
# transition is what enforces clause-10 termination: the newline is
# only consumed iff the in-progress value should continue past it.
# See issue #48.
trailer_val = (print | (nl when trailer_val_continues))+;

trailer_init = trailer_tok_breaking >mark @err(rewind) trailer_sep_breaking >set_current_footer_key @err(rewind) |
               trailer_tok >mark @err(rewind) trailer_sep >set_current_footer_key @err(rewind);

# Count newlines that can be part of the body or just trailer separators.
# Optionally match a trailer token followed by a trailer separator.
# In such case, continue looking for a trailer value.
# Otherwise, assume machine is in the body part.
trailer_beg := nl* $count_nl (trailer_init @complete_trailer_parsing)?;

# Match a trailer value.
# Then, ignoring newlines, continue trying to detect other trailers.
# set_footer is fired explicitly on the terminating boundary (a non-
# continuing newline OR EOF) instead of as a leaving action on
# trailer_val itself, because Ragel's `%`-style leaving actions fire
# on every loop-back transition inside `(...)+ ` — which would emit
# the value after every consumed continuation newline (issue #48).
trailer_end := (trailer_val >mark) $eof(set_footer)
               ((nl when !trailer_val_continues) @set_footer $count_nl @start_trailer_parsing)?;

# Match anything until two newlines (ie., a blank line).
# Then, try detect a footer looking for a trailer token.
body := (any >mark $err(append_body) %append_body %err(append_body_before_blank_line) when !blank_line_ahead)+ $err(start_trailer_parsing);

# Expect a blank line after the description.
# Try detect a footer looking for a trailer token.
remainder = blank_line @start_trailer_parsing;

main := minimal_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	scope? %to(check_early_exit)
	breaking? %to(check_early_exit)
	colon >err(err_colon) %to(check_early_exit)
	description
	remainder?;

conventional_types_main := conventional_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	scope? %to(check_early_exit)
	breaking? %to(check_early_exit)
	colon >err(err_colon) %to(check_early_exit)
	description
	remainder?;

falco_types_main := falco_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	scope? %to(check_early_exit)
	breaking? %to(check_early_exit)
	colon >err(err_colon) %to(check_early_exit)
	description
	remainder?;

free_form_types_main := free_form_types >eof(err_empty) >mark @err(err_type) %from(set_type) %to(check_early_exit)
	:> (scope? %to(check_early_exit))
	breaking? %to(check_early_exit)
	:>> colon >err(err_colon) %to(check_early_exit)
	description
	remainder?;

}%%

%% write data noerror noprefix;

type machine struct {
	data             []byte
	cs               int
	p, pe, eof       int
	pb               int
	err              error
	bestEffort       bool
	typeConfig       conventionalcommits.TypeConfig
	logger           *logrus.Logger
	currentFooterKey string
	countNewlines    int
	lastNewline      int
	// trailerBlockStart is the byte offset (within data) at which the
	// trailing footer trailer block begins, or -1 when no such block
	// exists. It is computed exactly once at the start of Parse() and
	// is consulted by shouldRedirectToBody to decide whether a given
	// >start_trailer_parsing entry should fall through to the body
	// machine instead. See issue #38.
	trailerBlockStart int
	// trailerLineStarts is the sorted list of byte offsets within data
	// at which a real trailer line begins (the first line of the block
	// plus every continuation that is itself trailer_init-shaped).
	// Empty when no trailer block is present. It is consulted by
	// trailerValueContinues to enforce spec clause-10 termination of
	// multi-line trailer values: a value continues across a newline
	// iff the next line is NOT registered here. See issue #48.
	trailerLineStarts []int
}

func (m *machine) text() []byte {
	return m.data[m.pb:m.p]
}

// shouldRedirectToBody is consulted by the start_trailer_parsing Ragel
// action to decide whether the next stretch of input should be parsed
// as body content rather than as the start of a footer trailer.
//
// The decision is anchored on a single pre-scan (findTrailerBlockStart)
// that is computed once per Parse() call and stored on the machine.
// When the FSM is currently positioned strictly before that pre-scanned
// trailer-block boundary, the upcoming line is body content (a
// trailer-shaped paragraph head, a fake-trailer line in the middle of
// a body paragraph, or trailing prose); otherwise it is — by
// construction of the pre-scan — the start of the real trailer block
// and the FSM should keep its historical behavior (enter trailer_beg).
//
// When the pre-scan found no trailer block at all (trailerBlockStart
// is -1) the predicate returns false: there is nothing to redirect to,
// trailer_beg can consume any trailing newlines on its own, and the
// existing rewind action handles the case where the FSM mistakenly
// ventured into trailer_beg with body content still ahead.
//
// The reason the predicate uses strict `<` (not `<=`) is that at
// `m.p == m.trailerBlockStart - 1` the FSM is already sitting on the
// final newline of the blank line that separates body from trailer; on
// the next byte (m.p+1 == trailerBlockStart) the trailer-token first
// character is the next thing the FSM will read. Letting trailer_beg
// take over from that byte is exactly what we want.
func (m *machine) shouldRedirectToBody() bool {
	// Find the position of the next non-newline byte starting from
	// m.p+1. At every callsite of start_trailer_parsing the FSM is
	// positioned such that m.p still belongs to the current production
	// (last byte of body consumed, or last byte of a footer value, or
	// last byte of the description). The relevant question is: does
	// the NEXT non-newline byte (start of the next line) live before
	// the pre-scanned trailer block?
	pos := m.p + 1
	for pos < m.pe && m.data[pos] == '\n' {
		pos++
	}
	if pos >= m.pe {
		// Only newlines/EOF remain. Nothing for body to do; let the
		// trailer_beg path absorb the trailing newlines on its own.
		return false
	}
	if m.trailerBlockStart < 0 {
		// Pre-scan found no trailing trailer block at all. The bytes
		// ahead must be body content (or input the FSM later rejects
		// as a malformed trailer; the existing rewind action handles
		// that case by demoting back to body).
		return true
	}

	// Redirect to body iff the next non-newline byte still lives
	// strictly before the start of the pre-scanned trailer block. When
	// pos == trailerBlockStart the FSM is positioned exactly on the
	// trailer token's first byte and the historical trailer_beg path
	// is what we want.
	return pos < m.trailerBlockStart
}

// trailerValueContinues is consulted by the trailer_val Ragel
// production to decide whether the in-progress multi-line trailer
// value should continue past the newline at m.p (clause-10
// termination predicate).
//
// The newline is consumed iff:
//   - the byte right after it is NOT the start of a real trailer line
//     per the pre-scan (otherwise the next line is the next trailer
//     and the value must terminate here), AND
//   - the next byte pair is not the start of a blank-line gap
//     (a blank line terminates the trailer block per the standard
//     FSM contract; nothing about clause-10 changes that).
//
// The first condition is the new clause-10 behavior. The second
// condition preserves pre-existing semantics. See issue #48.
func (m *machine) trailerValueContinues() bool {
	// Reject if we're at or past EOF: nothing to continue into.
	if m.p+1 >= m.pe {
		return false
	}
	// Reject if a blank line follows: the trailer block ends there.
	if m.p+2 < m.pe && m.data[m.p+1] == '\n' && m.data[m.p+2] == '\n' {
		return false
	}
	// Reject if the next line is itself a real trailer line per the
	// pre-scan. Binary search keeps this O(log n) per call.
	next := m.p + 1
	lo, hi := 0, len(m.trailerLineStarts)
	for lo < hi {
		mid := int(uint(lo+hi) >> 1)
		if m.trailerLineStarts[mid] < next {
			lo = mid + 1
		} else {
			hi = mid
		}
	}
	if lo < len(m.trailerLineStarts) && m.trailerLineStarts[lo] == next {
		return false
	}

	return true
}

func (m *machine) emitInfo(s string, args... interface{}) {
	if m.logger != nil {
		logEntry := logrus.NewEntry(m.logger)
		for i := 0; i < len(args); i = i + 2 {
			logEntry = m.logger.WithField(args[0].(string), args[1])
		}
		logEntry.Infoln(s)
	}
}

func (m *machine) emitDebug(s string, args... interface{}) {
	if m.logger != nil {
		logEntry := logrus.NewEntry(m.logger)
		for i := 0; i < len(args); i = i + 2 {
			logEntry = m.logger.WithField(args[0].(string), args[1])
		}
		logEntry.Debugln(s)
	}
}

func (m *machine) emitError(s string, args... interface{}) error {
	e := fmt.Errorf(s + ColumnPositionTemplate, args...)
	if m.logger != nil {
		m.logger.Errorln(e)
	}
	return e
}

func (m *machine) emitErrorWithoutCharacter(messageTemplate string) error {
	return m.emitError(messageTemplate, m.p)
}

func (m *machine) emitErrorOnCurrentCharacter(messageTemplate string) error {
	return m.emitError(messageTemplate, string(m.data[m.p]), m.p)
}

func (m *machine) emitErrorOnPreviousCharacter(messageTemplate string) error {
	return m.emitError(messageTemplate, string(m.data[m.p - 1]), m.p)
}

// NewMachine creates a new FSM able to parse Conventional Commits.
func NewMachine(options ...conventionalcommits.MachineOption) conventionalcommits.Machine {
	m := &machine{}

	for _, opt := range options {
		opt(m)
	}

	%% access m.;
	%% variable p m.p;
	%% variable pe m.pe;
	%% variable eof m.eof;
	%% variable data m.data;

	return m
}

// Parse parses the input byte array as a Conventional Commit message with no body neither footer.
//
// When a valid Conventional Commit message is given it outputs its structured representation.
// If the parsing detects an error it returns it with the position where the error occurred.
//
// It can also partially parse input messages returning a partially valid structured representation
// and the error that stopped the parsing.
func (m *machine) Parse(input []byte) (conventionalcommits.Message, error) {
	m.data = input
	m.p = 0
	m.pb = 0
	m.pe = len(input)
	m.eof = len(input)
	m.err = nil
	m.currentFooterKey = ""
	m.countNewlines = 0
	// Pre-compute the start offset of the trailing trailer block (or -1)
	// AND the sorted offsets of every real-trailer line within that
	// block, once per Parse() call. The FSM consults trailerBlockStart
	// (via shouldRedirectToBody) to keep body parsing going past any
	// trailer-shaped line that does not belong to the real trailer
	// block (issue #38), and trailerLineStarts (via
	// trailerValueContinues) to enforce spec clause-10 termination of
	// multi-line trailer values (issue #48).
	m.trailerBlockStart, m.trailerLineStarts = findTrailerBlock(input)
	output := &conventionalCommit{}
	output.footers = make(map[string][]string)
	output.typeconfig = m.typeConfig

	switch m.typeConfig {
	case conventionalcommits.TypesFreeForm:
		m.cs = en_free_form_types_main
		break
	case conventionalcommits.TypesFalco:
		m.cs = en_falco_types_main
		break
	case conventionalcommits.TypesConventional:
		m.cs = en_conventional_types_main
		break
	case conventionalcommits.TypesMinimal:
		fallthrough
	default:
		%% write init;
		break
	}
	%% write exec;

	if m.cs < first_final {
		if m.bestEffort && output.minimal() {
			// An error occurred but partial parsing is on and partial message is minimally valid
			return output.export(), m.err
		}
		// Defensive: a few FSM error transitions (e.g. a trailing CR
		// after an otherwise-complete trailer) leave m.cs < first_final
		// without ever assigning m.err. Returning (nil, nil) at this
		// point silently violates Parse's contract and nil-derefs the
		// canonical caller pattern:
		//
		//   m, err := p.Parse(in)
		//   if err != nil { return err }
		//   if !m.Ok() { ... }            // <- NPE here on (nil, nil)
		//
		// Synthesize an early-exit error from the current cursor so
		// strict-mode Parse always returns either a message or an
		// error, never both nil. The cursor may be at EOF (m.p == m.pe)
		// when this fires, so report against the previous byte when one
		// exists and against an empty character otherwise. See PR #47
		// second-pass review (§3) and TestParseNeverReturnsNilNil.
		if m.err == nil {
			if m.p > 0 && m.p-1 < len(m.data) {
				m.err = m.emitErrorOnPreviousCharacter(ErrEarly)
			} else {
				m.err = m.emitError(ErrEarly, "", m.p)
			}
		}
		return nil, m.err
	}

	return output.export(), nil
}

// WithBestEffort enables best effort mode.
func (m *machine) WithBestEffort() {
	m.bestEffort = true
}

// HasBestEffort tells whether the receiving machine has best effort mode on or off.
func (m *machine) HasBestEffort() bool {
	return m.bestEffort
}

// WithTypes tells the parser which commit message types to consider.
func (m *machine) WithTypes(t conventionalcommits.TypeConfig) {
	m.typeConfig = t
}

// WithLogger tells the parser which logger to use.
func (m *machine) WithLogger(l *logrus.Logger) {
	m.logger = l
}