local env = require("go:env")
local V = require("go:vars")
local shell = require("go:shell")

local M = {}

function M.get(name)
    local vt = name:sub(1,1)
    if vt == "%" then
        return tostring(V[name:sub(2)] or "")
    else
        local nNum = tonumber(name)
        if nNum and nNum > 0 then
            return shell.getArg(nNum) or ""
        end
        return env[name] or ""
    end
end

function M.set(name, value)
    local vt = name:sub(1,1)
    if vt == "%" then
        V[name:sub(2)] = value
    else
        env[name] = value
    end
end

return M
