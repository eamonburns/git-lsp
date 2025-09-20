-- Minimal init.lua to test git-lsp
--
-- Run with:
-- ```
-- nvim -u min_init.lua
-- ```
--
-- NOTE: Requires NeoVim 0.11+

if not vim.uv.fs_stat("./build/git-lsp") then
	vim.notify("Unable to execute ./build/git-lsp", vim.log.levels.ERROR)
end

vim.lsp.config["git-lsp"] = {
	cmd = { "./build/git-lsp" },

	filetypes = { "lua" },

	root_markers = { ".git" },
}

vim.lsp.enable("git-lsp")

vim.o.winborder = "rounded"
