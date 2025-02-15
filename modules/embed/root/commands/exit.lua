local shell = require("go:shell")

local M = {}

function M.run(_, code)
    shell.exit(tonumber(code))
end

return M
