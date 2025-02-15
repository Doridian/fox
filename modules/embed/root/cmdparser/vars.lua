local fs = require("fox.fs")
local Env = require("fox.env")

local M = {}

function M.get(varType, name)
    if varType == "$" then
        return Env[name] or ""
    elseif varType == "%" then
        return tostring(_G[name])
    end
end

function M.set(varType, name, value)
    if varType == "$" then
        Env[name] = value
    elseif varType == "%" then
        _G[name] = value
    end
end

-- return true to indicate that glob processing mode should be enabled
function M.interpolate(str, escapeGlobs)
    local i = 1
    local varStart, varEnd, varTmp, varType

    local hasGlobs = false
    local retStrEscaped = nil
    local retStr = ""

    while i <= #str do
        varStart = str:find("[$%%]", i)
        if not varStart then
            break
        end
        varType = str:sub(varStart, varStart)

        varTmp = str:sub(varStart + 1, varStart + 1)
        if varTmp == "{"  then
            varStart = varStart + 1
            varEnd = str:find("}", varStart + 1, true)
            if not varEnd then
                return nil, "Unclosed variable ${}"
            end
        else
            varEnd = str:find("[^%w_]", varStart + 1)
        end
        if not varEnd then
            varEnd = #str + 1
        end
        varTmp = str:sub(varStart + 1, varEnd - 1)
        varTmp = M.get(varType, varTmp)

        retStr = retStr .. varTmp
        if hasGlobs then
            retStrEscaped = retStrEscaped .. fs.globEscape(varTmp)
        elseif escapeGlobs and fs.hasGlob(varTmp) then
            hasGlobs = true
            retStrEscaped = retStr .. fs.globEscape(varTmp)
        end

        i = varEnd + 1
    end

    retStr = retStr .. str:sub(i)
    if hasGlobs then
        retStrEscaped = retStrEscaped .. str:sub(i)
    end
    return retStr, retStrEscaped, hasGlobs
end

return M
