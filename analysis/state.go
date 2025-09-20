package analysis

import (
	"fmt"
	"strings"

	"github.com/eamonburns/git-lsp/lsp"
)

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: make(map[string]string)}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	diagnostcs := []lsp.Diagnostic{}

	vsCode := "VS Code"
	neoVim := "NeoVim"

	for row, line := range strings.Split(text, "\n") {
		if strings.Contains(line, vsCode) {
			idx := strings.Index(line, vsCode)
			diagnostcs = append(diagnostcs, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len(vsCode)),
				Severity: 1,
				Source:   "Common sense",
				Message:  "Please make sure we use good language in this video",
			})
		}

		if strings.Contains(line, neoVim) {
			idx := strings.Index(line, neoVim)
			diagnostcs = append(diagnostcs, lsp.Diagnostic{
				Range:    LineRange(row, idx, idx+len(neoVim)),
				Severity: 3,
				Source:   "Common sense",
				Message:  "Great choice",
			})
		}
	}

	return diagnostcs
}

func (self *State) OpenDocument(uri string, text string) []lsp.Diagnostic {
	self.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (self *State) UpdateDocument(uri string, text string) []lsp.Diagnostic {
	self.Documents[uri] = text

	return getDiagnosticsForFile(text)
}

func (self *State) Hover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	document := self.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("# Document attributes\n\n- URI: %s\n- Characters: %d", uri, len(document)),
		},
	}
}

func (self *State) TextDocumentCompletion(id int, uri string, position lsp.Position) lsp.CompletionResponse {
	// TODO: Get completions
	items := []lsp.CompletionItem{
		{
			Label:         "NeoVim (BTW)",
			Detail:        "Very cool editor",
			Documentation: "Fun to watch in videos :)",
			Kind:          lsp.CompletionItemKindText,
		},
		{
			Label:         "Position",
			Detail:        fmt.Sprintf("line: %d, character: %d", position.Line, position.Character),
			Documentation: "Current position in the file",
			Kind:          lsp.CompletionItemKindText,
		},
	}

	return lsp.CompletionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: items,
	}
}

func LineRange(line int, start int, end int) lsp.Range {
	return lsp.Range{
		Start: lsp.Position{
			Line:      line,
			Character: start,
		},
		End: lsp.Position{
			Line:      line,
			Character: end,
		},
	}
}
