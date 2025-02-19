local shell = require("embed:parsers.shell")

local M = {}

function M.run(_, script)
    shell.runLine(script)
end

return M
