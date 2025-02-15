local fs = require("fox.fs")

shell.commands = {}

local function interpolateVars(str)
    return str
end

function shell.parsers.cmd(cmd, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(cmd, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    local i = 1
    local args = {}
    local curArg, nextControlIdx, nextControl, quoteEndIdx
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

        if container ~= "'" then
            sub = interpolateVars(sub)
        end

        if curArg.isGlob then
            if container then
                subEscaped = fs.globEscape(sub)
            else
                subEscaped = sub
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
    return parsed
end
