local env = require("go:env")

local M = {}

function M.run(_, _, varSet)
    local eqPos = varSet:find("=", 1, true)
    if not eqPos then
        error("missing variable value")
        return 1
    end
    local varKey = varSet:sub(1, eqPos - 1)
    local varVal = varSet:sub(eqPos + 1)
    _G[varKey] = varVal
    return 0
end

return M
