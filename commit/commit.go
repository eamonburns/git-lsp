package commit

import (
	"fmt"
	"strings"

	"github.com/eamonburns/git-lsp/internal/helper"
	"github.com/eamonburns/git-lsp/lsp"
)

// TODO: Better severities

// Conventional Commits Specification: https://www.conventionalcommits.org/en/v1.0.0/#specification

type Commit struct {
	// e.g. feat, fix, docs
	Type string

	Scope string

	// Description of the breaking change (if any)
	// If the breaking change is specified by a "!" in the type/scope
	// prefix, and there is no BREAKING CHANGE footer, then this will
	// be the same as Description
	BreakingChange string

	Description string

	Body string

	Footers map[string]string
}

func Parse(text string) (Commit, []Diagnostic) {
	diagnostics := []Diagnostic{}
	header, rest, _ := strings.Cut(text, "\n")

	commit := Commit{}

	typeScope, description, foundTypeScope := strings.Cut(header, ": ")
	if foundTypeScope {
		commit.Description = description

		// Check for breaking change "!"
		if idx := strings.LastIndex(typeScope, "!"); idx != -1 {
			// Remove "!"
			typeScope = typeScope[:idx] // Check to make sure

			// "13. If included in the type/scope prefix, breaking changes MUST be indicated by a ! immediately before the :. If ! is used, BREAKING CHANGE: MAY be omitted from the footer section, and the commit description SHALL be used to describe the breaking change."
			// NOTE: I am interpreting the above requirement to mean that if the BREAKING CHANGE footer is included, that is used instead of the description
			commit.BreakingChange = description
		}

		// Extract scope if present
		if lParIdx := strings.Index(typeScope, "("); lParIdx != -1 {
			commit.Type = typeScope[:lParIdx]

			rParIdx := strings.Index(typeScope, ")")

			if rParIdx == -1 {
				commit.Scope = typeScope[lParIdx+1:]

				diagnostics = append(diagnostics, Diagnostic{
					Range: helper.LineRange(0, lParIdx, lParIdx),
					Type:  UnmatchedLeftParenError,
				})
			} else if rParIdx < len(typeScope)-1 {
				// The right parentheses is not the last character of typeScope

				commit.Scope = typeScope[lParIdx+1 : rParIdx]

				diagnostics = append(diagnostics, Diagnostic{
					Range: helper.LineRange(0, rParIdx+1, len(typeScope)),
					Type:  ExtraCharactersAfterScopeError,
					Args:  []string{typeScope[rParIdx+1:]},
				})
			} else {
				commit.Scope = typeScope[lParIdx+1 : rParIdx]
			}
		} else if idx := strings.Index(typeScope, ")"); idx != -1 {
			// There wasn't a '(', but there was a ')'

			diagnostics = append(diagnostics, Diagnostic{
				Range: helper.LineRange(0, idx, idx),
				Type:  UnmatchedRightParenError,
			})
			commit.Type = typeScope[:idx]
		} else {
			commit.Type = typeScope
		}
	} else {
		diagnostics = append(diagnostics, Diagnostic{
			Range: helper.LineRange(0, 0, len(header)),
			Type:  NoTypeScopeError,
		})
		commit.Description = typeScope // Header line wasn't split, so typeScope is the whole line, which we will use as the description
	}

	// TODO: Parse body/footers
	// NOTE: Ignore all lines starting with #
	_ = rest

	if foundTypeScope && commit.Type == "" {
		diagnostics = append(diagnostics, Diagnostic{
			Range: helper.LineRange(0, 0, 0),
			Type:  EmptyTypeError,
		})
	}
	if idx := strings.Index(typeScope, "("); idx != -1 && foundTypeScope && commit.Scope == "" {
		diagnostics = append(diagnostics, Diagnostic{
			Range: helper.LineRange(0, len(commit.Type), len(commit.Type)+2),
			Type:  EmptyScopeError,
		})
	}

	return commit, diagnostics
}

type DiagnosticType int

// Diagnostic error/warning types
const (
	// There was no type/scope in the header line (e.g. "description")
	NoTypeScopeError DiagnosticType = iota
	// There was a left parentheses in the type/scope, but no matching right parentheses (e.g. "type(scope: description")
	UnmatchedLeftParenError
	// There was a right parentheses in the type/scope, but no matching left parentheses (e.g. "typescope): description")
	UnmatchedRightParenError
	// There were extra characters after the scope (e.g. "type(scope)bla: description")
	// Args: 0 = characters
	ExtraCharactersAfterScopeError
	// The type in the type/scope was empty (e.g. "(scope): description")
	EmptyTypeError
	// The scope in the type/scope was empty (e.g. "type(): description")
	EmptyScopeError
)

type Diagnostic struct {
	Range lsp.Range
	Type  DiagnosticType
	Args  []string
}

func (self Diagnostic) ToLspDiagnostic() lsp.Diagnostic {
	var message string
	var severity int = 1

	switch self.Type {
	case NoTypeScopeError:
		message = "No type/scope in header line"
	case UnmatchedLeftParenError:
		message = "Unmatched '('"
	case UnmatchedRightParenError:
		message = "Unmatched ')'"
	case ExtraCharactersAfterScopeError:
		message = fmt.Sprintf("Extra characters after scope: '%s'", self.Args[0])
	case EmptyTypeError:
		message = "Empty type"
	case EmptyScopeError:
		message = "Empty scope"
	default:
		message = "Unknown error"
	}

	return lsp.Diagnostic{
		Range:    self.Range,
		Severity: severity,
		Source:   "git-lsp",
		Message:  message,
	}
}
