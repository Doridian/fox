local shell = require("fox.shell")

local M = {}

function M.run(code)
    shell.exit(code)
end
M.name = "exit"

return M
