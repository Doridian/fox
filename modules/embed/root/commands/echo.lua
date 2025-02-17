local M = {}

function M.run(ctx, _, ...)
    ctx.stdout:print(...)
    return 0
end

M.canLua = true
M.mustLua = true

return M
