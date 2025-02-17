local os = require("go:os")

local M = {}

function M.runDirect(_, dir)
    os.chdir(dir)
    return 0
end

return M
