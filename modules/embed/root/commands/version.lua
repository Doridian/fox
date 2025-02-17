local info = require("go:info")

local M = {}

function M.run(ctx)
    ctx.stdout:print("version: " .. info.version)
    ctx.stdout:print("commit: " .. info.gitrev)
    return 0
end

M.mustLua = true

return M
