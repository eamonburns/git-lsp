package commit

import (
	"testing"

	"github.com/eamonburns/git-lsp/internal/helper"
	"github.com/eamonburns/git-lsp/lsp"
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

	assert.ElementsMatch(t, []lsp.Diagnostic{
		{
			Range:    helper.LineRange(0, 0, len(commitMsg)-1),
			Severity: 1,
			Source:   "git-lsp",
			Message:  "No type/scope in header line",
		},
	}, diagnostics)
	assert.Equal(t, Commit{
		Description: commitMsg,
	}, commit)

	// TODO: No type: "(scope): description"
	// TODO: Empty scope: "type(): description"
	// TODO: No description: "type(scope):"
	// TODO: Extra characters after scope: "type(scope)bla: description"
}
