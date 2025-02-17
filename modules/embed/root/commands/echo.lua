local M = {}

function M.run(ctx, _, ...)
    ctx.stdout:write(table.concat({...}, " ") .. "\n")
    return 0
end

return M
