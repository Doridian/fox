local os = require("go:os")

local M = {}

function M.run(ctx)
    ctx.stdout:print(os.getwd())
    return 0
end

M.canLua = true
M.mustLua = true

return M
