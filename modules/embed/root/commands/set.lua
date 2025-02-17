local env = require("go:env")

local M = {}

function M.runDirect(_, varSet)
    local eqPos = varSet:find("=", 1, true)
    if not eqPos then
        _G[varSet] = env[varSet] or _G[varSet]
        return 0
    end
    local varKey = varSet:sub(1, eqPos - 1)
    local varVal = varSet:sub(eqPos + 1)
    _G[varKey] = varVal
    return 0
end

return M
