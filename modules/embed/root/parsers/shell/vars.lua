local env = require("go:env")

local M = {}

function M.get(varType, name)
    if varType == "$" then
        return env[name] or ""
    elseif varType == "%" then
        return tostring(_G[name] or "")
    end
end

function M.set(varType, name, value)
    if varType == "$" then
        env[name] = value
    elseif varType == "%" then
        _G[name] = value
    end
end

return M
