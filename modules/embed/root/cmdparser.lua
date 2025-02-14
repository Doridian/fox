
shell.commands = {}

function shell.parsers.cmd(cmd, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(cmd, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    local i = 1
    local args = {}
    local curArg
    local nextControlIdx, nextControl, quoteEndIdx
    local function bufArg(quoted)
        if not curArg then
            curArg = {
                value = "",
                quoted = false,
            }
        end
        if quoted then
            curArg.quoted = true
        end

        if nextControlIdx then
            curArg.value = curArg.value .. parsed:sub(i, nextControlIdx - 1)
            i = nextControlIdx + 1
        else
            curArg.value = curArg.value .. parsed:sub(i)
            i = #parsed + 1
        end
    end
    local function pushArg(quoted)
        bufArg(quoted)
        if curArg.value ~= "" then
            table.insert(args, curArg)
            curArg = nil
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
            bufArg(true)
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
