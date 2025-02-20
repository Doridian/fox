local shell = require("embed:parsers.shell")
local goshell = require("go:shell")

local M = {}

function M.run(_, script)
    -- TOOD: Change this to argOffset param in shell module
    goshell.popArgs(2)
    shell.runLine(script)
end

return M
