package analysis

import "github.com/eamonburns/git-lsp/lsp"

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: make(map[string]string)}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	return []lsp.Diagnostic{
		{
			Range:    LineRange(0, 0, 4),
			Severity: 1,
			Source:   "my brain",
			Message:  "There is a thing",
		},
	}
}

func (self *State) OpenDocument(uri string, text string) []lsp.Diagnostic {
	self.Documents[uri] = text

	return getDiagnosticsForFile(text)
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
