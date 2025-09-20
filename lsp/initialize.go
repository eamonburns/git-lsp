package lsp

// Type definitions for "initialize" request
// https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/#initialize

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	// There is a ton more that could go here
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}
