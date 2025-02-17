local fs = require("go:fs")
local Env = require("go:env")

local M = {}

function M.get(varType, name)
    if varType == "$" then
        return Env[name] or ""
    elseif varType == "%" then
        return _G[name] or ""
    end
end

function M.set(varType, name, value)
    if varType == "$" then
        Env[name] = value
    elseif varType == "%" then
        _G[name] = value
    end
end

return M
