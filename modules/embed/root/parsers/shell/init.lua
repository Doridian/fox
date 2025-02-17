local os = require("go:os")
local tokenizer = require("embed:parsers.shell.tokenizer")
local splitter = require("embed:parsers.shell.splitter")
local gocmd = require("go:cmd")
local fs = require("go:fs")
local cmdHandler = require("embed:commandHandler")
local pipe = require("go:pipe")

local exe = os.executable()

local M = {}

local function cmdRun(cmd)
    if cmd._runPre then
        pcall(cmd._runPre)
    end

    local ctx = {
        stdin = cmd._stdin,
        stdout = cmd._stdout,
        stderr = cmd._stderr,
    }
    local dummyCtx = {}
    if cmd._ref_stdin then
        table.insert(dummyCtx, ctx.stdin)
        ctx.stdin = cmd["_"..cmd._ref_stdin]
    end
    if cmd._ref_stdout then
        table.insert(dummyCtx, ctx.stdout)
        ctx.stdout = cmd["_"..cmd._ref_stdout]
    end
    if cmd._ref_stderr then
        table.insert(dummyCtx, ctx.stderr)
        ctx.stderr = cmd["_"..cmd._ref_stderr]
    end
    local ok, exitCode = pcall(cmd.run, ctx, cmd.args)
    if not ok then
        (ctx.stderr or pipe.stderr):write(exitCode)
        exitCode = 1
    end
    cmdHandler.closeCtx(ctx)
    cmdHandler.closeCtx(dummyCtx)

    if cmd._runPost then
        pcall(cmd._runPost)
    end

    return exitCode or 0
end

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
        local fh, err = fs.open(redir.name, fMode)
        if not fh then
            error(err)
        end
        if not cmd.gocmd then
            cmd["_"..name] = fh
        else
            cmd.gocmd[name](cmd.gocmd, fh)
        end
    elseif redir.type == splitter.RedirTypeRefer then
        cmd["_ref_"..name] = redir.ref
        cmd["_ref_" .. redir.ref] = "null"

        if cmd.gocmd then
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
                cmd._stdin = redir.cmd.gocmd:stdoutPipe()
                cmd._runPre = function()
                    redir.cmd.gocmd:start()
                end
                cmd._runPost = function()
                    redir.cmd.gocmd:wait()
                end
            else
                -- mostly ignore it, sb -> sb doesn't do anything but order
                local subPipe = pipe.new()
                redir.cmd._stdout = subPipe
                cmd._stdin = subPipe
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
        pipe.stderr:write("shell.tokenizer error: " .. err .. "\n")
        return ""
    end
    local cmds, err = splitter.run(tokens)
    if not cmds then
        pipe.stderr:write("shell.splitter error: " .. err .. "\n")
        return ""
    end

    local rootCmds = {}

    for _, cmd in pairs(cmds) do
        cmd._null = pipe.null
        cmd._stdin = pipe.stdin
        cmd._stdout = pipe.stdout
        cmd._stderr = pipe.stderr

        local cmdObj, _ = cmdHandler.get(cmd.args[1])
        if cmdObj then
            if cmdObj.forbidInline then
                table.insert(cmd.args, 1, exe)
                table.insert(cmd.args, 2, "-c")
            else
                cmd.run = function(ctx, subargs)
                    return cmdObj.run(ctx, table.unpack(subargs))
                end
            end
        end
        if not cmd.run then
            cmd.gocmd = gocmd.new(cmd.args)
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
