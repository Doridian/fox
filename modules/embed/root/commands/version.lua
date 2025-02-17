local info = require("go:info")

local M = {}

function M.run(ctx)
    ctx.stdout:write("version: " .. info.version .. "\ncommit: " .. info.gitrev .. "\n")
    return 0
end

return M
