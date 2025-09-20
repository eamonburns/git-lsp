package lsp

type Request struct {
	RPC    string `json:"jsonrpc"`
	ID     int    `json:"id"`
	Method string `json:"method"`

	// We will want to specify the types of parameters in all the Request types
	// Params ...
}
