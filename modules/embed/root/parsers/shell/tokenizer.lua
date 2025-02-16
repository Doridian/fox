local fs = require("go:fs")
local interpolate = require("embed:parsers.shell.interpolate")

local M = {}

M.ArgTypeString = 1
M.ArgTypeOp = 2

function M.run(parsed)
    local i = 1
    local tokens = {}
    local curToken, nextControlIdx, nextControl, controlEndIdx, foundGlobs
    local function bufToken(container)
        if not curToken then
            curToken = {
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
            sub, subEscaped, foundGlobs = interpolate.run(sub, not container)
            if not sub then
                -- subEscaped will be the error message
                return nil, "Parse error: " .. tostring(subEscaped)
            end

            if (not container) and foundGlobs then
                curToken.isGlob = true
            end
        end

        if curToken.isGlob then
            if container then
                subEscaped = fs.globEscape(sub)
            end
            curToken.bufEscaped = curToken.bufEscaped .. subEscaped
        elseif (not container) and fs.hasGlob(sub) then
            curToken.isGlob = true
            -- We can only get here if nothing previously could be a glob
            -- so we can just escape everything in buf lazily here
            curToken.bufEscaped = fs.globEscape(curToken.buf) .. sub
        end

        curToken.buf = curToken.buf .. sub
    end
    local function manualPushToken()
        if #curToken.buf > 0 then
            local arg = curToken
            curToken = nil
            if arg.isGlob then
                local matches = fs.glob(arg.bufEscaped)
                if #matches > 0 then
                    for _, match in pairs(matches) do
                        table.insert(tokens, {
                            val = match,
                            type = M.ArgTypeString,
                        })
                    end
                    return
                end
            end
            table.insert(tokens, {
                val = arg.buf,
                type = M.ArgTypeString,
            })
        end
    end
    local function pushToken(container)
        bufToken(container)
        manualPushToken()
    end
    while i <= #parsed do
        nextControlIdx = parsed:find("[ \n\t\"';&|><!]", i)
        if not nextControlIdx then
            pushToken()
            break
        end

        nextControl = parsed:sub(nextControlIdx, nextControlIdx)
        if nextControl == "'" or nextControl == '"' then
            controlEndIdx = parsed:find(nextControl, nextControlIdx + 1)
            if not controlEndIdx then
                return parsed, true, nextControl .. "> "
            end

            bufToken()
            nextControlIdx = controlEndIdx
            bufToken(nextControl)
        elseif nextControl == "\n" or nextControl == "\r" or nextControl == "\t" or nextControl == " " then
            pushToken()
        else
            bufToken()
            if (nextControl ~= "<" and nextControl ~= ">") or not (curToken and tonumber(curToken.buf)) then
                manualPushToken()
            end

            controlEndIdx = nextControlIdx
            while parsed:sub(controlEndIdx + 1, controlEndIdx + 1) == nextControl do
                controlEndIdx = controlEndIdx + 1
            end

            table.insert(tokens, {
                pre = curToken and curToken.buf,
                val = nextControl,
                len = controlEndIdx - nextControlIdx + 1,
                raw = parsed:sub(nextControlIdx, controlEndIdx),
                type = M.ArgTypeOp,
            })
            curToken = nil
            i = controlEndIdx + 1
        end
    end

    return tokens
end

return M
