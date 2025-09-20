package analysis

import "github.com/eamonburns/git-lsp/lsp"

type State struct {
	Documents map[string]string
}

func NewState() State {
	return State{Documents: make(map[string]string)}
}

func getDiagnosticsForFile(text string) []lsp.Diagnostic {
	return []lsp.Diagnostic{}
}

func (self *State) OpenDocument(uri string, text string) []lsp.Diagnostic {
	self.Documents[uri] = text

	return getDiagnosticsForFile(text)
}
