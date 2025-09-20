package helper

import "github.com/eamonburns/git-lsp/lsp"

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
