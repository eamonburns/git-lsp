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
			Range: helper.LineRange(0, 0, len(commitMsg)-1),
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

	// TODO: No type: "(scope): description"
	// TODO: Empty scope: "type(): description"
	// TODO: No description: "type(scope):"
	// TODO: Unmatched '(': "type(scope: description"
	// TODO: Unmatched ')': "typescope): description"
}
