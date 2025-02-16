local cmd = require("go:cmd")
local fs = require("go:fs")
local shell = require("go:shell")
local os = require("go:os")
local interpolate = require("embed:multiparser.shell.interpolate")

local exe = os.executable()

--[[
    TODO:
    - Need to do command splitting at parse time before variable/glob interp
        - | || & && ;
    - Need todo stdio redirections
]]

local M = {}

function M.run(str, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(str, lineNo)
    if (not parsed) or parsed == true or parsed == "" or parsed == "\n" then
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
            sub, subEscaped, foundGlobs = interpolate.run(sub, not container)
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

    if #args < 1 then
        return
    end

    if shell.hasCommand(args[1]) then
        table.insert(args, 1, exe)
        table.insert(args, 2, "-c")
    end
    local sCmd = cmd.new(args)
    sCmd:errorPropagation(false)
    -- TODO: Assign multi-command here, assign stdio here
    return function()
        sCmd:run()
    end
end

return M
