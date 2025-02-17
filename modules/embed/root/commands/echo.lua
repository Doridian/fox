local M = {}

function M.run(ctx, _, ...)
    ctx.stdout:print(...)
    return 0
end

M.mustLua = true

return M
