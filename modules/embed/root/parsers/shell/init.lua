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

local function setGocmdStdio(cmd, name)
    local redir = cmd[name]
    if not redir then
        return
    end

    if redir.type == splitter.RedirTypeFile then
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
        local refObj
        if redir.ref == "stdout" then
            refObj = cmd.gocmd:stdoutPipe()
        elseif redir.ref == "stderr" then
            refObj = cmd.gocmd:stderrPipe()
        elseif redir.ref == "stdin" then
            refObj = cmd.gocmd:stdinPipe()
        else
            error("invalid refer type: " .. tostring(redir.ref))
        end
        cmd.gocmd[name](cmd.gocmd, refObj)
    elseif redir.type == splitter.RedirTypeCmd then
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
        cmd.gocmd = gocmd.new(args)

        setGocmdStdio(cmd, "stdin")
        setGocmdStdio(cmd, "stdout")
        setGocmdStdio(cmd, "stderr")

        if wait and not cmd.stdin then
            cmd.gocmd:stdin(shell.stdin, false)
        end
        cmd.gocmd:start()
        table.insert(cmds, 1, cmd)
        if cmd.stdin and cmd.stdin.type == splitter.RedirTypeCmd then
            cmd = cmd.stdin.cmd
        else
            cmd = nil
        end
    end

    if not wait then
        return 0
    end

    local exitOK = true
    for _, c in pairs(cmds) do
        local exitCodeS = c.gocmd:wait()
        local exitOKSub = exitCodeS == 0
        if c.invert then
            exitOKSub = not exitOKSub
        end
        if (not exitOKSub) and c ~= rootCmd then
            exitOK = false
            if errorOnPipeFail then
                shell.stderr:print("piped command " .. tostring(c and c.gocmd) .. " exited with code " .. exitCodeS)
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

    local tokens, err = tokenizer.run(parsed)
    if not tokens then
        shell.stderr:print("shell.tokenizer error", err)
        return ""
    end

    local cmds, err = splitter.run(tokens)
    if not cmds then
        shell.stderr:print("shell.splitter error", err)
        return ""
    end

    local rootCmds = {}

    for _, cmd in pairs(cmds) do
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
        for _, cmd in pairs(rootCmds) do
            if cmd.background then
                startGoCmdDeep(cmd, false)
            else
                if not skipNext then
                    local okSub = startGoCmdDeep(cmd, true)
                    exitCode = cmd.gocmd:wait()
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
                    shell.stderr:print("invalid chainToNext: " .. tostring(cmd.chainToNext))
                    return
                end
            end
        end

        if errorOnFail and not exitSuccess then
            shell.stderr:print("command " .. tostring(exitCmd and exitCmd.gocmd) .. " exited with code " .. exitCode)
        end
    end
end

return M
