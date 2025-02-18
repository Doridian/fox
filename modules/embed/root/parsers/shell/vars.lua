local env = require("go:env")
local V = require("go:vars")

local M = {}

function M.get(varType, name)
    if varType == "$" then
        return env[name] or ""
    elseif varType == "%" then
        return tostring(V[name] or "")
    end
end

function M.set(varType, name, value)
    if varType == "$" then
        env[name] = value
    elseif varType == "%" then
        V[name] = value
    end
end

return M
