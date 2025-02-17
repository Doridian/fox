local shell = require("go:shell")

local M = {}

function M.run(_, _, code)
    shell.exit(code)
    return code
end

M.canLua = true
M.mustLua = true

return M
