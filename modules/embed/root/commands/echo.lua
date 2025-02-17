local os = require("go:os")

local M = {}

function M.runDirect(_, ...)
    return 0, table.concat({...}, " ") .. "\n"
end

return M

