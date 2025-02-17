local M = {}

function M.run(ctx, _, code)
---@diagnostic disable-next-line: deprecated
    local func = loadstring(code, "commands.eval")
    if not func then
        return 1
    end

    local ret = tostring(func() or "")
    if ret == "" then
        return 0
    end
    ctx.stdout:print(ret)
    return 0
end

M.canLua = true
M.mustLua = true

return M
