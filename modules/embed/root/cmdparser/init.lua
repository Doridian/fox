local cmd = require("fox.cmd")
local fs = require("fox.fs")
local Env = require("fox.env")
local shell = require("fox.shell")
local vars = require("fox.embed.cmdparser.vars")
local cmds = require("fox.embed.cmdparser.cmds")

--[[
    TODO:
    - Implement a stack of stdios for redirection to work
    - Implement a way for Lua's print to respect stdout changes (likely implement custom print)
    - Need to :wait() for cmd processes that produce input for lua stuff before running it
        - Likely change runCmd -> loadCmd
]]

local M = {}

M.cmds = cmds

local function loadCmd(args)
    local cmdName = args[1]
    if not cmdName then
        return
    end

    local luaCmd = cmds.get(cmdName)
    if luaCmd then
        table.remove(args, 1)
        return function()
            luaCmd(unpack(args))
        end
    end
end

function M.parser(str, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(str, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    local i = 1
    local args = {}
    local curArg, nextControlIdx, nextControl, quoteEndIdx, foundGlobs
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

        subEscaped = sub
        if container ~= "'" then
            sub, subEscaped, foundGlobs = vars.interpolate(sub, not container)
            if not sub then
                -- subEscaped will be the error message
                print("Parse error: " .. tostring(subEscaped))
                return ""
            end

            if (not container) and foundGlobs then
                curArg.isGlob = true
            end
        end

        if curArg.isGlob then
            if container then
                subEscaped = fs.globEscape(sub)
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

    local cmdFunc = loadCmd(args)
    if not cmdFunc then
        local sCmd = cmd.new(args)
        cmdFunc = function()
            return sCmd:run()
        end
    end
    local exitCode = tonumber(cmdFunc() or 0)
    return "_G._LAST_EXIT_CODE = " .. tostring(exitCode)
end

return M
