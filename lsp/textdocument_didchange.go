package lsp

type DidChangeTextDocumentNotification struct {
	Notification
	Params DidChangeTextDocumentParams `json:"params"`
}

type DidChangeTextDocumentParams struct {
	TextDocument   VersionTextDocumentIdentifier    `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// An event describing a change to a text document
// If only `Text` is provided, it is considered to be the full content of the document.
type TextDocumentContentChangeEvent struct {
	// New text of the whole document
	Text string `json:"text"`
}
