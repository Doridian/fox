local os = require("go:os")
local tokenizer = require("embed:parsers.shell.tokenizer")
local splitter = require("embed:parsers.shell.splitter")
local gocmd = require("go:cmd")
local fs = require("go:fs")
local cmdHandler = require("embed:commandHandler")

local exe = os.executable()

local M = {}

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
        cmd.gocmd[name](cmd.gocmd, fh)
    elseif redir.type == splitter.RedirTypeCmd then
        if name == "stdin" then
            cmd.gocmd:stdin(redir.cmd.gocmd:stdoutPipe())
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
        local args = cmd.args
        if cmdHandler.has(cmd.args[1]) then
            table.insert(args, 1, exe)
            table.insert(args, 2, "-c")
        end
        cmd.gocmd = gocmd.new(args)
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
        for _, cmd in pairs(rootCmds) do
            if cmd.background then
                cmd.gocmd:start()
            else
                if not skipNext then
                    local exitCode = cmd.gocmd:run()
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
