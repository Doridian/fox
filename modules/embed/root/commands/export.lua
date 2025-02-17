local env = require("go:env")

local M = {}

function M.runDirect(_, envSet)
    local eqPos = envSet:find("=", 1, true)
    if not eqPos then
        error("missing variable value")
        return 1
    end
    local envKey = envSet:sub(1, eqPos - 1)
    local envVal = envSet:sub(eqPos + 1)
    env[envKey] = envVal
    return 0
end

return M
