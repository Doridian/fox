local M = {}

M.parsers = {}

function M.loadParser(name, func)
    M.parsers[name] = func
end
local function loadIntegratedParser(name)
    M.loadParser(name, require("embed:multiparser." .. name).run)
end
loadIntegratedParser("lua")
loadIntegratedParser("shell")

function M.run(cmd, lineNo, prev)
    if cmd:sub(1, 1) == "!" then
        local newLine = cmd:find("\n", 1, true)
        local cmdPrefix = cmd:sub(2, newLine - 1)
        if M.parsers[cmdPrefix] then
            return M.parsers[cmdPrefix](cmd:sub(newLine + 1), lineNo, prev)
        end

        print("Unknown parser " .. cmdPrefix)
        return ""
    end

    local defParser = M.parsers.default
    if defParser then
        return defParser(cmd, lineNo, prev)
    end
    return false
end

return M
