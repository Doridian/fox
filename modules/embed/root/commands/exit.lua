local shell = require("fox.shell")

local M = {}

function M.run(code)
    shell.exit(tonumber(code))
end

return M
