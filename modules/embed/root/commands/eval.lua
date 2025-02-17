local M = {}

function M.runDirect(_, code)
---@diagnostic disable-next-line: deprecated
    local func = loadstring(code, "commands.eval")
    if not func then
        return 1
    end

    local ret = tostring(func() or "")
    if ret == "" then
        return 0
    end
    return 0, ret .. "\n"
end

return M
