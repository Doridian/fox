local shell = require("fox.shell")

local M = {}

function M.run(_, code)
    shell.exit(tonumber(code))
end

return M
