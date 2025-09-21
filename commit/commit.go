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
	// If the breaking change is specified by a "!" in the type/scope prefix, then this will be the same as Description
	BreakingChange string

	Description string

	Body string

	Footers map[string]string
}

func Parse(text string) (Commit, []Diagnostic) {
	diagnostics := []Diagnostic{}
	header, rest, found := strings.Cut(text, "\n")

	commit := Commit{}

	typeScope, description, found := strings.Cut(header, ": ")
	if found {
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
			Range: helper.LineRange(0, 0, len(header)-1),
			Type:  NoTypeScopeError,
		})
		commit.Description = typeScope // Header line wasn't split, so typeScope is the whole line, which we will use as the description
	}

	// TODO: Parse body/footers
	_ = rest

	return commit, diagnostics
}

type DiagnosticType int

// Diagnostic error/warning types
const (
	// There was no type/scope in the header line
	NoTypeScopeError DiagnosticType = iota
	// There was a right parentheses in the type/scope, but no matching left parentheses
	UnmatchedRightParenError
	// There was a left parentheses in the type/scope, but no matching right parentheses
	UnmatchedLeftParenError
	// There were extra characters after the scope
	// Args: 0 = characters
	ExtraCharactersAfterScopeError
)

type Diagnostic struct {
	Range lsp.Range
	Type  DiagnosticType
	Args  []string
}

func (self Diagnostic) ToLspDiagnostic() lsp.Diagnostic {
	var message string
	var severity int

	switch self.Type {
	case NoTypeScopeError:
		message = "No type/scope in header line"
		severity = 1
	case UnmatchedLeftParenError:
		message = "Unmatched '('"
		severity = 1
	case UnmatchedRightParenError:
		message = "Unmatched ')'"
		severity = 1
	case ExtraCharactersAfterScopeError:
		message = fmt.Sprintf("Extra characters after scope: '%s'", self.Args[0])
		severity = 1
	default:
		message = "Unknown error"
		severity = 1
	}

	return lsp.Diagnostic{
		Range:    self.Range,
		Severity: severity,
		Source:   "git-lsp",
		Message:  message,
	}
}
