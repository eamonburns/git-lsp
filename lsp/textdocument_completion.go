package lsp

type CompletionRequest struct {
	Request
	Params CompletionParams `json:"params"`
}

type CompletionParams struct {
	TextDocumentPositionParams
}

type CompletionResponse struct {
	Response
	Result []CompletionItem `json:"result"`
}

type CompletionItem struct {
	Label         string             `json:"label"`
	Detail        string             `json:"detail"`
	Documentation string             `json:"documentation"`
	Kind          CompletionItemKind `json:"kind"`
}

type CompletionItemKind int

const (
	CompletionItemKindText CompletionItemKind = iota + 1
	CompletionItemKindMethod
	CompletionItemKindFunction
	CompletionItemKindConstructor
	CompletionItemKindField
	CompletionItemKindVariable
	CompletionItemKindClass
	CompletionItemKindInterface
	CompletionItemKindModule
	CompletionItemKindProperty
	CompletionItemKindUnit
	CompletionItemKindValue
	CompletionItemKindEnum
	CompletionItemKindKeyword
	CompletionItemKindSnippet
	CompletionItemKindColor
	CompletionItemKindFile
	CompletionItemKindReference
	CompletionItemKindFolder
	CompletionItemKindEnumMember
	CompletionItemKindConstant
	CompletionItemKindStruct
	CompletionItemKindEvent
	CompletionItemKindOperator
	CompletionItemKindTypeParameter
)
