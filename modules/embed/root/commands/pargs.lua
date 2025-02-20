local shell = require("go:shell")

local M = {}

function M.run(...)
    local args = {...}
    for i, v in ipairs(args) do
        print(i, v)
    end
end

return M
