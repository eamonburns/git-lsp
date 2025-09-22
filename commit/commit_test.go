package commit

import (
	"testing"

	"github.com/eamonburns/git-lsp/internal/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	commit, diagnostics := Parse("type(scope)!: description")
	require.Empty(t, diagnostics)
	assert.Equal(t, Commit{
		Type:           "type",
		Scope:          "scope",
		BreakingChange: "description",
		Description:    "description",
	}, commit)

	commit, diagnostics = Parse("type: description")
	require.Empty(t, diagnostics)
	assert.Equal(t, Commit{
		Type:           "type",
		Scope:          "",
		BreakingChange: "",
		Description:    "description",
	}, commit)

	commitMsg := "description"
	commit, diagnostics = Parse(commitMsg)
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 0, len(commitMsg)),
			Type:  NoTypeScopeError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Description: commitMsg,
	}, commit)

	commit, diagnostics = Parse("type(scope)bla: description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 11, 14),
			Type:  ExtraCharactersAfterScopeError,
			Args:  []string{"bla"},
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "type",
		Scope:       "scope",
		Description: "description",
	}, commit)

	commit, diagnostics = Parse("type(scope: description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 4, 4),
			Type:  UnmatchedLeftParenError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "type",
		Scope:       "scope",
		Description: "description",
	}, commit)

	commit, diagnostics = Parse("typescope): description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 9, 9),
			Type:  UnmatchedRightParenError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "typescope",
		Scope:       "",
		Description: "description",
	}, commit)

	commit, diagnostics = Parse("(scope): description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 0, 0),
			Type:  EmptyTypeError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "",
		Scope:       "scope",
		Description: "description",
	}, commit)

	commit, diagnostics = Parse("type(): description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 4, 6),
			Type:  EmptyScopeError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "type",
		Scope:       "",
		Description: "description",
	}, commit)

	commit, diagnostics = Parse("(): description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 0, 2),
			Type:  EmptyScopeError,
		},
		{
			Range: helper.LineRange(0, 0, 0),
			Type:  EmptyTypeError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "",
		Scope:       "",
		Description: "description",
	}, commit)

	commit, diagnostics = Parse("type(scope):")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 12, 12),
			Type:  EmptyDescriptionError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "type",
		Scope:       "scope",
		Description: "",
	}, commit)

	commit, diagnostics = Parse("type(scope):description")
	assert.ElementsMatch(t, []Diagnostic{
		{
			Range: helper.LineRange(0, 12, 12),
			Type:  NoSpaceBeforeDescriptionError,
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Type:        "type",
		Scope:       "scope",
		Description: "description",
	}, commit)
}

func TestParseFooter(t *testing.T) {
	_, _, ok := ParseFooter("not a footer")
	require.False(t, ok)

	footer, value, ok := ParseFooter("hash-footer # value")
	require.True(t, ok)
	assert.Equal(t, "hash-footer", footer)
	assert.Equal(t, "value", value)

	footer, value, ok = ParseFooter("colon-footer: value")
	require.True(t, ok)
	assert.Equal(t, "colon-footer", footer)
	assert.Equal(t, "value", value)

	footer, value, ok = ParseFooter("whitespace footer: value")
	require.False(t, ok)

	footer, value, ok = ParseFooter("colon-footer: # value")
	require.True(t, ok)
	assert.Equal(t, "colon-footer", footer)
	assert.Equal(t, "# value", value)

	footer, value, ok = ParseFooter("hash-footer #: value")
	require.True(t, ok)
	assert.Equal(t, "hash-footer", footer)
	assert.Equal(t, ": value", value)

	// Breaking change whitespace exception
	footer, value, ok = ParseFooter("BREAKING CHANGE: value")
	require.True(t, ok)
	assert.Equal(t, "BREAKING CHANGE", footer)
	assert.Equal(t, "value", value)

	footer, value, ok = ParseFooter("BREAKING CHANGE # value")
	require.False(t, ok)

	footer, value, ok = ParseFooter("BREAKING CHANGE # : value")
	require.False(t, ok)
}
