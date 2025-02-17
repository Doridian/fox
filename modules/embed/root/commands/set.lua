local M = {}

function M.run(ctx, _, varSet)
    local eqPos = varSet:find("=", 1, true)
    local varKey, varVal
    if not eqPos then
        varKey = varSet
        varVal = ctx.stdin:readToEnd()
    else
        varKey = varSet:sub(1, eqPos - 1)
        varVal = varSet:sub(eqPos + 1)
    end
    _G[varKey] = varVal
    return 0
end

M.mustLua = true

return M
