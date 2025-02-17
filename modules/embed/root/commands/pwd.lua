local os = require("go:os")

local M = {}

function M.run(ctx)
    ctx.stdout:write(os.getwd() .. "\n")
    return 0
end

return M
