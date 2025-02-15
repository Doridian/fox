local fs = require("fox.fs")
local Env = require("fox.env")

shell.commands = {}

local function fetchVar(varType, varName)
    if varType == "$" then
        return Env[varName] or ""
    elseif varType == "%" then
        return tostring(_G[varName])
    end
end

-- return true to indicate that glob processing mode should be enabled
local function interpolateVars(str, escapeGlobs)
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
                error("Unclosed variable ${}")
            end
        else
            varEnd = str:find("[^%w_]", varStart + 1)
        end
        if not varEnd then
            varEnd = #str + 1
        end
        varTmp = str:sub(varStart + 1, varEnd - 1)
        varTmp = fetchVar(varType, varTmp)

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

function shell.parsers.cmd(cmd, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(cmd, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    local i = 1
    local args = {}
    local curArg, nextControlIdx, nextControl, quoteEndIdx, foundGlobs
    local function bufArg(container)
        if not curArg then
            curArg = {
                buf = "",
                bufEscaped = "",
                isGlob = false,
            }
        end

        local sub, subEscaped
        if nextControlIdx then
            sub = parsed:sub(i, nextControlIdx - 1)
            i = nextControlIdx + 1
        else
            sub = parsed:sub(i)
            i = #parsed + 1
        end

        subEscaped = sub
        if container ~= "'" then
            sub, subEscaped, foundGlobs = interpolateVars(sub, not container)
            if (not container) and foundGlobs then
                curArg.isGlob = true
            end
        end

        if curArg.isGlob then
            if container then
                subEscaped = fs.globEscape(sub)
            end
            curArg.bufEscaped = curArg.bufEscaped .. subEscaped
        elseif (not container) and fs.hasGlob(sub) then
            curArg.isGlob = true
            -- We can only get here if nothing previously could be a glob
            -- so we can just escape everything in buf lazily here
            curArg.bufEscaped = fs.globEscape(curArg.buf) .. sub
        end

        curArg.buf = curArg.buf .. sub
    end
    local function pushArg(container)
        bufArg(container)
        if #curArg.buf > 0 then
            local arg = curArg
            curArg = nil
            if arg.isGlob then
                local matches = fs.glob(arg.bufEscaped)
                if #matches > 0 then
                    for _, match in pairs(matches) do
                        table.insert(args, match)
                    end
                    return
                end
            end
            table.insert(args, arg.buf)
        end
    end
    while i <= #parsed do
        nextControlIdx = parsed:find("[ \n\t\"']", i)
        if not nextControlIdx then
            pushArg()
            break
        end

        nextControl = parsed:sub(nextControlIdx, nextControlIdx)
        if nextControl == "'" or nextControl == '"' then
            quoteEndIdx = parsed:find(nextControl, nextControlIdx + 1)
            if not quoteEndIdx then
                return true, nextControl .. "> "
            end

            bufArg()
            nextControlIdx = quoteEndIdx
            bufArg(nextControl)
        else
            pushArg()
        end
    end

    for k, v in pairs(args) do
        print("ARG", k, v)
    end

    -- TODO: Parse CLI-like language
    return ""
end
