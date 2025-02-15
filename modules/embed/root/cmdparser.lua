
shell.commands = {}

local function interpolateVars(str)
    return str
end

local function globStar(arg)
end

function shell.parsers.cmd(cmd, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(cmd, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    local i = 1
    local args = {}
    local buf = {}
    local nextControlIdx, nextControl, quoteEndIdx
    local function bufArg(container)
        local sub
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
        table.insert(buf, sub)
    end
    local function pushArg(container)
        bufArg(container)
        if #buf > 0 then
            table.insert(args, buf)
            buf = {}
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
        print("ARG", k, v.value, v.quoted)
    end

    -- TODO: Parse CLI-like language
    return parsed
end
