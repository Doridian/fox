local os = require("go:os")

local M = {}

function M.run(_, dir)
    os.chdir(dir)
    return 0
end

return M
