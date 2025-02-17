local os = require("go:os")

local M = {}

function M.runDirect()
    return 0, os.getwd() .. "\n"
end

return M
