local os = require("go:os")
local tokenizer = require("embed:parsers.shell.tokenizer")
local splitter = require("embed:parsers.shell.splitter")
local gocmd = require("go:cmd")
local fs = require("go:fs")
local cmdHandler = require("embed:commandHandler")
local shell = require("go:shell")
local pipe = require("go:pipe")

local exe = os.executable()

local M = {}

local function setGocmdStdio(cmd, name)
    local redir = cmd[name]
    if not redir then
        return
    end

    if not cmd.gocmd then
        error("cannot redirect superbuiltin: " .. tostring(cmd.args[1]))
    end

    if redir.type == splitter.RedirTypeFile then
        local fh, err = fs.open(redir.name, redir.append and "a" or "w")
        if not fh then
            error(err)
        end
        if not cmd.gocmd then
            if name == "stdout" then
                cmd._stdout = fh
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
            elseif redir.cmd.gocmd then
                cmd._stdout = redir.cmd.gocmd:stdoutPipe()
                cmd._runPre = function()
                    return redir.cmd.gocmd:run()
                end
            else
                -- mostly ignore it, sb -> sb doesn't do anything but order
                cmd._runPre = function()
                    return redir.cmd.run(redir.cmd.args)
                end
            end
        else
            error("cannot pipe cmd into stdout or stderr")
        end
    else
        error("invalid redir type: " .. tostring(redir.type))
    end
end

-- TODO: Make Lua cmds stdout a reader that when first read from runs the code

local superBuiltins = {}
function superBuiltins.cd(args)
    os.chdir(args[2])
    return 0
end
function superBuiltins.exit(args)
    shell.exit(args[2])
    return 0
end
function superBuiltins.pwd(_)
    return 0, os.getwd()
end
function superBuiltins.echo(args)
    return 0, table.concat(args, " ", 2)
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
        local args = cmd.args
        if superBuiltins[args[1]] then
            cmd.run = superBuiltins[args[1]]
        elseif cmdHandler.has(cmd.args[1]) then
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
        local stdout
        local exitCode
        for _, cmd in pairs(rootCmds) do
            if cmd.background then
                if cmd.run then
                    if cmd._runPre then
                        _, _ = cmd._runPre()
                    end
                    _, stdout = cmd.run(cmd.args)
                    if cmd._stdout then
                        cmd._stdout:write(stdout)
                        cmd._stdout:close()
                    else
                        print(stdout)
                    end
                else
                    cmd.gocmd:start()
                end
            else
                if not skipNext then
                    if cmd.run then
                        if cmd._runPre then
                            _, _ = cmd._runPre()
                        end
                        exitCode, stdout = cmd.run(cmd.args)
                        if cmd._stdout then
                            cmd._stdout:write(stdout)
                            cmd._stdout:close()
                        else
                            print(stdout)
                        end
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
