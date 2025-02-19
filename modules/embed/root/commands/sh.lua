local goShell = require("go:shell")
local shell = require("embed:parsers.shell")

local M = {}

function M.run(_, script)
    goShell.popArgs(2)
    shell.runLine(script)
end

return M
