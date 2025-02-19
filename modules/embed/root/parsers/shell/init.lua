local shell = require("go:shell")
local tokenizer = require("embed:parsers.shell.tokenizer")
local splitter = require("embed:parsers.shell.splitter")
local gocmd = require("go:cmd")
local fs = require("go:fs")

-- TODO?: Implement "(echo A && echo B) | grep A" type subshells
-- TODO?: Implement \ escaping

local errorOnFail = true
local errorOnPipeFail = true

local M = {}

local GETTER_FUNCS = {
    stdin = "getStdin",
    stdout = "getStdout",
    stderr = "getStderr",
}

local PIPE_FUNCS = {
    stdin = "stdinPip",
    stdout = "stdoutPipe",
    stderr = "stderrPipe",
}

local function setGocmdStdio(cmd, name, phase)
    local redir = cmd[name]
    if not redir then
        if phase ~= 2 or name == "stdin" then
            return
        end
        if not cmd.gocmd[GETTER_FUNCS[name]](cmd.gocmd) then
            cmd.gocmd[name](cmd.gocmd, shell[name], false)
        end
        return
    end

    if redir.type == splitter.RedirTypeFile then
        if phase ~= 1 then
            return
        end

        local fMode
        if name == "stdin" then
            fMode = "r"
        elseif redir.append then
            fMode = "a"
        else
            fMode = "w"
        end
        local fh, err = fs.open(tokenizer.oneStringVal(redir.name), fMode)
        if not fh then
            error(err)
        end

        cmd.gocmd[name](cmd.gocmd, fh)
    elseif redir.type == splitter.RedirTypeRefer then
        if phase ~= 2 then
            return
        end

        local refObj = cmd.gocmd[GETTER_FUNCS[redir.ref]](cmd.gocmd)
        cmd.gocmd[name](cmd.gocmd, refObj, false)
    elseif redir.type == splitter.RedirTypeCmd then
        if phase ~= 1 then
            return
        end

        if name ~= "stdin" then
            error("cannot pipe cmd into stdout or stderr")
        end

        -- attach cmd's stdin to redir.cmd's stdout
        cmd.gocmd:stdin(redir.cmd.gocmd:stdoutPipe())
    else
        error("invalid redir type: " .. tostring(redir.type))
    end
end

local function startGoCmdDeep(rootCmd, wait)
    local cmds = {}
    local cmd = rootCmd
    while cmd do
        local args = {}
        for _, arg in pairs(cmd.args) do
            tokenizer.stringVals(arg, args)
        end
        cmd.gocmd:args(args)

        setGocmdStdio(cmd, "stdin", 1)
        setGocmdStdio(cmd, "stdout", 1)
        setGocmdStdio(cmd, "stderr", 1)

        if wait and not cmd.stdin then
            cmd.gocmd:stdin(shell.stdin, false)
        end

        table.insert(cmds, 1, cmd)
        if cmd.stdin and cmd.stdin.type == splitter.RedirTypeCmd then
            cmd = cmd.stdin.cmd
        else
            cmd = nil
        end
    end

    for _, cmd in pairs(cmds) do
        setGocmdStdio(cmd, "stdin", 2)
        setGocmdStdio(cmd, "stdout", 2)
        setGocmdStdio(cmd, "stderr", 2)

        pcall(cmd.gocmd.start, cmd.gocmd)
    end

    if not wait then
        return 0
    end

    local exitOK = true
    local ok, err
    for _, c in pairs(cmds) do
        local ok, exitCodeS = pcall(c.gocmd.wait, c.gocmd)
        if not ok then
            err = exitCodeS
            exitCodeS = 1
        end
        local exitOKSub = exitCodeS == 0
        if c.invert then
            exitOKSub = not exitOKSub
        end
        if (not exitOKSub) and c ~= rootCmd then
            exitOK = false
            if errorOnPipeFail then
                shell.stderr:print("piped command " .. tostring(c.gocmd) .. " " .. (err or ("exited with code " .. exitCodeS)))
            end
        end
    end
    return exitOK
end

function M.run(strAdd, lineNo, prev)
    local parsed = (prev or "") .. strAdd .. "\n"

    if strAdd:sub(#strAdd, #strAdd) == "\\" then
        return parsed:sub(1, #parsed - 2) .. "\n", true
    end

    return M.runLine(parsed)
end

function M.runLine(parsed)
    local tokens, err = tokenizer.run(parsed)
    if not tokens then
        error("shell.tokenizer error " .. tostring(err))
    end

    local cmds, err = splitter.run(tokens)
    if not cmds then
        error("shell.splitter error " .. tostring(err))
    end

    local rootCmds = {}

    for _, cmd in pairs(cmds) do
        cmd.gocmd = gocmd.new()
        rootCmds[cmd] = cmd
    end

    -- Do this after rootCmds is prefilled
    for _, cmd in pairs(cmds) do
        if cmd.stdin and cmd.stdin.type == splitter.RedirTypeCmd then
            rootCmds[cmd.stdin.cmd] = nil
        end
    end

    -- for _, cmd in pairs(cmds) do
    --     shell.stderr:print(cmd.args[1], cmd.run and "lua" or "cmd")
    -- end

    return function()
        local skipNext = false
        local exitSuccess = true
        local exitCode = 0
        local exitCmd = nil
        local err, ok
        for _, cmd in pairs(rootCmds) do
            if cmd.background then
                startGoCmdDeep(cmd, false)
            else
                if not skipNext then
                    local okSub = startGoCmdDeep(cmd, true)
                    ok, exitCode = pcall(cmd.gocmd.wait, cmd.gocmd)
                    if not ok then
                        err = exitCode
                        exitCode = 1
                    end
                    exitCmd = cmd
                    exitSuccess = exitCode == 0
                    if cmd.invert then
                        exitSuccess = not exitSuccess
                    end

                    if errorOnPipeFail and not okSub then
                        break
                    end
                else
                    skipNext = false
                end

                if cmd.chainToNext == "&&" then
                    if not exitSuccess then
                        skipNext = true
                    end
                elseif cmd.chainToNext == "||" then
                    if exitSuccess then
                        skipNext = true
                    end
                elseif cmd.chainToNext then
                    error("invalid chainToNext: " .. tostring(cmd.chainToNext))
                end
            end
        end

        if errorOnFail and not exitSuccess then
            error("command " .. tostring(exitCmd and exitCmd.gocmd) .. " " .. (err or ("exited with code " .. exitCode)))
        end
    end
end

return M
