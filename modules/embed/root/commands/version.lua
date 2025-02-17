local info = require("go:info")

local M = {}

function M.runDirect()
    return 0, "version: " .. info.version .. "\ncommit: " .. info.gitrev .. "\n"
end

return M
