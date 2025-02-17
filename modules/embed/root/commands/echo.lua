local os = require("go:os")

local M = {}

function M.run(ctx, ...)
    ctx.stdout:write(table.concat({...}, " ") .. "\n")
    return 0
end

return M

