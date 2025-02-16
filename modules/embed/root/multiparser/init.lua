local shell = require("go:shell")

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

function M.run(cmdAdd, lineNo, prev)
    if not prev then
        prev = {
            prev = nil,
        }
        local cmdPrefix = cmdAdd:sub(1, 1)
        if cmdPrefix == "!" then
            cmdAdd = cmdAdd:sub(2)
            prev.parser = M.parsers[cmdAdd]

            if not prev.parser then
                print("Unknown parser " .. cmdAdd)
                return ""
            end

            return prev, true
        elseif cmdPrefix == "=" or cmdAdd:sub(1, 2) == "--" then
            prev.parser = shell.defaultShellParser
        elseif cmdPrefix == "/" then
            prev.parser = shell.defaultShellParser
            cmdAdd = cmdAdd:sub(2)
        else
            prev.parser = M.parsers.default
        end
    end

    if not prev.parser then
        return false
    end

    local state, needMore, promptOverride = prev.parser(cmdAdd, lineNo, prev.prev)
    if needMore then
        prev.prev = state
        return prev, true, promptOverride
    end
    return state, false, promptOverride
end

return M
