
shell.commands = {}

function shell.parsers.cmd(cmd, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(cmd, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    -- TODO: Parse CLI-like language
    return parsed
end
