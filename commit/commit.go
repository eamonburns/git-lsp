package commit

import (
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

func Parse(text string) (Commit, []lsp.Diagnostic) {
	diagnostics := []lsp.Diagnostic{}
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
		if idx := strings.Index(typeScope, "("); idx != -1 {
			commit.Type = typeScope[:idx]
			scope := typeScope[idx+1:]

			idx := strings.Index(scope, ")")

			if idx == -1 {
				commit.Scope = scope
			} else if idx < len(scope)-1 {
				commit.Scope = scope[:idx]

				diagnostics = append(diagnostics, lsp.Diagnostic{
					Range:    helper.LineRange(0, idx+1, len(typeScope)-1), // FIXME: Invalid range
					Severity: 1,
					Source:   "git-lsp",
					Message:  "Extra characters after end of scope: '%s'",
				})
			} else {
				commit.Scope = scope[:idx]
			}
		} else if idx := strings.Index(typeScope, ")"); idx != -1 {
			diagnostics = append(diagnostics, lsp.Diagnostic{
				Range:    helper.LineRange(0, idx, idx),
				Severity: 1,
				Source:   "git-lsp",
				Message:  "Unmatched ')'",
			})
			commit.Type = typeScope[:idx]
		} else {
			commit.Type = typeScope
		}
	} else {
		diagnostics = append(diagnostics, lsp.Diagnostic{
			Range:    helper.LineRange(0, 0, len(header)-1),
			Severity: 1,
			Source:   "git-lsp",
			Message:  "No type/scope in header line",
		})
		commit.Description = typeScope // Header line wasn't split, so typeScope is the whole line, which we will use as the description
	}

	// TODO: Parse body/footers
	_ = rest

	return commit, diagnostics
}
