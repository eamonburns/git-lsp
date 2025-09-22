package commit

import (
	"testing"

	"github.com/eamonburns/git-lsp/internal/helper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	commit, diagnostics := Parse("feat(thing)!: cool")
	require.Empty(t, diagnostics)

	assert.Equal(t, Commit{
		Type:           "feat",
		Scope:          "thing",
		BreakingChange: "cool",
		Description:    "cool",
	}, commit)

	commit, diagnostics = Parse("hi: ollo")
	require.Empty(t, diagnostics)

	assert.Equal(t, Commit{
		Type:           "hi",
		Scope:          "",
		BreakingChange: "",
		Description:    "ollo",
	}, commit)

	commitMsg := "no type or scope"
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

	// TODO: Empty type: "(scope): description"
	// TODO: Empty scope: "type(): description"
	// TODO: No description: "type(scope):"
}
