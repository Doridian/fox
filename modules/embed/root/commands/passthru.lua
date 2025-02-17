local os = require("go:os")

local M = {}

function M.run(ctx)
    while true do
        local data = ctx.stdin:read(1024)
        if not data then
            break
        end
        ctx.stdout:write(data)
    end
    return 0
end

return M

