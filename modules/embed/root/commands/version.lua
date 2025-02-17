local info = require("go:info")

local M = {}

function M.run()
    print("version: " .. info.version)
    print("commit: " .. info.gitrev)
    return 0
end

return M
