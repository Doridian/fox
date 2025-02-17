local os = require("go:os")
local tokenizer = require("embed:parsers.shell.tokenizer")
local splitter = require("embed:parsers.shell.splitter")
local gocmd = require("go:cmd")
local fs = require("go:fs")
local cmdHandler = require("embed:commandHandler")
local pipe = require("go:pipe")

local exe = os.executable()

local M = {}

-- TODO: Implement merge runDirect + run
-- TODO: Implement ctx as first cmd arg such that { stderr = PIPE, stdout = PIPE, stdin = PIPE, name = arg0 } can be passed

local function cmdRun(cmd)
    if cmd._runPre then
        pcall(cmd._runPre)
    end

    local ctx = {
        stdin = cmd._stdin,
        stdout = cmd._stdout,
        stderr = cmd._stderr,
    }
    local ok, exitCode = pcall(cmd.run, ctx, cmd.args)
    if not ok then
        (ctx.stderr or pipe.stderr):write(exitCode)
        cmdHandler.closeCtx(ctx)
        exitCode = 1
    end

    return exitCode or 0
end

local function setGocmdStdio(cmd, name)
    local redir = cmd[name]
    if not redir then
        return
    end

    if redir.type == splitter.RedirTypeFile then
        local fh, err = fs.open(redir.name, redir.append and "a" or "w")
        if not fh then
            error(err)
        end
        if not cmd.gocmd then
            if name == "stdout" then
                cmd._stdout = fh
            elseif name == "stderr" then
                cmd._stderr = fh
            end
        else
            cmd.gocmd[name](cmd.gocmd, fh)
        end
    elseif redir.type == splitter.RedirTypeCmd then
        if name == "stdin" then
            if cmd.gocmd then
                if redir.cmd.gocmd then
                    cmd.gocmd:stdin(redir.cmd.gocmd:stdoutPipe())
                    return
                end
                redir.cmd._stdout = cmd.gocmd:stdinPipe()
                cmd.gocmd:addPreReq(function()
                    cmdRun(redir.cmd)
                end)
            elseif redir.cmd.gocmd then
                redir.cmd.gocmd:stdout(pipe.null)
                cmd._runPre = function()
                    redir.cmd.gocmd:run()
                end
            else
                -- mostly ignore it, sb -> sb doesn't do anything but order
                redir.cmd._stdout = pipe.null
                cmd._runPre = function()
                    cmdRun(redir.cmd)
                end
            end
        else
            error("cannot pipe cmd into stdout or stderr")
        end
    else
        error("invalid redir type: " .. tostring(redir.type))
    end
end

function M.run(strAdd, lineNo, prev)
    local parsed = (prev or "") .. strAdd .. "\n"

    if strAdd:sub(#strAdd, #strAdd) == "\\" then
        return parsed:sub(1, #parsed - 2) .. "\n", true
    end

    local tokens, err = tokenizer.run(parsed)
    if not tokens then
        print("shell.tokenizer error: " .. err)
        return ""
    end
    local cmds, err = splitter.run(tokens)
    if not cmds then
        print("shell.splitter error: " .. err)
        return ""
    end

    local rootCmds = {}

    for _, cmd in pairs(cmds) do
        local hasCmd, hasDirect = cmdHandler.has(cmd.args[1])
        local args = cmd.args
        if hasDirect then
            cmd.run = function(subargs)
                return cmdHandler.run(args[1], subargs, true)
            end
        elseif hasCmd then
            table.insert(args, 1, exe)
            table.insert(args, 2, "-c")
        end
        if not cmd.run then
            cmd.gocmd = gocmd.new(args)
        end
        rootCmds[cmd] = cmd
    end

    -- Do this after so all gocmd structures are for sure filled
    for _, cmd in pairs(cmds) do
        if cmd.stdin and cmd.stdin.type == splitter.RedirTypeCmd then
            rootCmds[cmd.stdin.cmd] = nil
        end

        setGocmdStdio(cmd, "stdin")
        setGocmdStdio(cmd, "stdout")
        setGocmdStdio(cmd, "stderr")
    end

    return function()
        local skipNext = false
        local exitSuccess = true
        local exitCode
        for _, cmd in pairs(rootCmds) do
            if cmd.background then
                if cmd.run then
                    cmdRun(cmd)
                else
                    cmd.gocmd:start()
                end
            else
                if not skipNext then
                    if cmd.run then
                        exitCode = cmdRun(cmd)
                    else
                        exitCode = cmd.gocmd:run()
                    end
                    exitSuccess = exitCode == 0
                    if cmd.invert then
                        exitSuccess = not exitSuccess
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
    end
end

return M
