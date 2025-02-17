local env = require("go:env")

local M = {}

function M.run(ctx, _, envSet)
    local eqPos = envSet:find("=", 1, true)
    local envKey, envVal
    if not eqPos then
        envKey = envSet
        envVal = ctx.stdin:readToEnd()
    else
        envKey = envSet:sub(1, eqPos - 1)
        envVal = envSet:sub(eqPos + 1)
    end
    env[envKey] = envVal
    return 0
end

M.canLua = true
M.mustLua = true

return M
