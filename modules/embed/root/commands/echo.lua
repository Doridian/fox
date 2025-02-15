local shell = require("fox.shell")

local M = {}

function M.run(...)
    print(...)
end
M.name = "echo"

return M
