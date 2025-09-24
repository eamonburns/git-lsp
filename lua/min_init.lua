-- Minimal init.lua to test git-lsp
--
-- Run with:
-- ```
-- nvim -u lua/min_init.lua
-- ```
--
-- NOTE: Requires NeoVim 0.11+

---------- Get executable path ----------

-- Directory containing this file
local lua_dir = vim.fs.dirname(vim.fs.abspath(debug.getinfo(1, "S").short_src))
-- Directory containing built git-lsp: normalize("lua/../build")
local build_dir = vim.fs.normalize(vim.fs.joinpath(lua_dir, "..", "build"))
-- Name of executable
local exe_name = "git-lsp"
if vim.fn.has("win32") == 1 and vim.fn.has("wsl") == 0 then
	exe_name = exe_name .. ".exe"
end
local exe_path = vim.fs.joinpath(build_dir, exe_name)
if not vim.uv.fs_stat(exe_path) then
	vim.notify("Unable to execute '" .. exe_path .. "'", vim.log.levels.ERROR)
end

---------- Configure LSP ----------

vim.lsp.config["git-lsp"] = {
	cmd = { exe_path },

	filetypes = { "gitcommit" },

	root_markers = { ".git" },
}
vim.lsp.enable("git-lsp")

---------- Configure Neovim ----------

vim.o.winborder = "rounded"
