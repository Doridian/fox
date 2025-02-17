local env = require("go:env")

local M = {}

function M.runDirect(_, envSet)
    local eqPos = envSet:find("=", 1, true)
    if not eqPos then
        env[envSet] = _G[envSet] or env[envSet]
        return 0
    end
    local envKey = envSet:sub(1, eqPos - 1)
    local envVal = envSet:sub(eqPos + 1)
    env[envKey] = envVal
    return 0
end

return M
