local duration = require("go:duration")

local M = {}

function M.run(_, _, durStr)
    duration.parse(durStr):sleepFor()
    return 0
end

return M
