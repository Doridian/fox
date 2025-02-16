local cmd = require("go:cmd")
local fs = require("go:fs")
local shell = require("go:shell")
local os = require("go:os")
local interpolate = require("embed:parsers.shell.interpolate")
local cmdHandler = require("embed:commandHandler")

local exe = os.executable()

--[[
    TODO:
    - Need to do command splitting at parse time before variable/glob interp
        - | || & && ;
    - Need todo stdio redirections
]]

local M = {}

local ArgTypeString = 1
local ArgTypeOp = 2

function M.run(strAdd, lineNo, prev)
    local parsed = (prev or "") .. strAdd .. "\n"

    if strAdd:sub(#strAdd, #strAdd) == "\\" then
        return parsed:sub(1, #parsed - 2) .. "\n", true
    end

    local i = 1
    local tokens = {}
    local curToken, nextControlIdx, nextControl, controlEndIdx, foundGlobs
    local function bufToken(container)
        if not curToken then
            curToken = {
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
                curToken.isGlob = true
            end
        end

        if curToken.isGlob then
            if container then
                subEscaped = fs.globEscape(sub)
            end
            curToken.bufEscaped = curToken.bufEscaped .. subEscaped
        elseif (not container) and fs.hasGlob(sub) then
            curToken.isGlob = true
            -- We can only get here if nothing previously could be a glob
            -- so we can just escape everything in buf lazily here
            curToken.bufEscaped = fs.globEscape(curToken.buf) .. sub
        end

        curToken.buf = curToken.buf .. sub
    end
    local function manualPushToken()
        if #curToken.buf > 0 then
            local arg = curToken
            curToken = nil
            if arg.isGlob then
                local matches = fs.glob(arg.bufEscaped)
                if #matches > 0 then
                    for _, match in pairs(matches) do
                        table.insert(tokens, {
                            val = match,
                            type = ArgTypeString,
                        })
                    end
                    return
                end
            end
            table.insert(tokens, {
                val = arg.buf,
                type = ArgTypeString,
            })
        end
    end
    local function pushToken(container)
        bufToken(container)
        manualPushToken()
    end
    while i <= #parsed do
        nextControlIdx = parsed:find("[ \n\t\"';&|><!]", i)
        if not nextControlIdx then
            pushToken()
            break
        end

        nextControl = parsed:sub(nextControlIdx, nextControlIdx)
        if nextControl == "'" or nextControl == '"' then
            controlEndIdx = parsed:find(nextControl, nextControlIdx + 1)
            if not controlEndIdx then
                return parsed, true, nextControl .. "> "
            end

            bufToken()
            nextControlIdx = controlEndIdx
            bufToken(nextControl)
        elseif nextControl == "\n" or nextControl == "\r" or nextControl == "\t" or nextControl == " " then
            pushToken()
        else
            bufToken()
            if (nextControl ~= "<" and nextControl ~= ">") or not (curToken and tonumber(curToken.buf)) then
                manualPushToken()
            end

            controlEndIdx = nextControlIdx
            while parsed:sub(controlEndIdx + 1, controlEndIdx + 1) == nextControl do
                controlEndIdx = controlEndIdx + 1
            end

            table.insert(tokens, {
                pre = curToken and curToken.buf,
                val = nextControl,
                len = controlEndIdx - nextControlIdx + 1,
                raw = parsed:sub(nextControlIdx, controlEndIdx),
                type = ArgTypeOp,
            })
            curToken = nil
            i = controlEndIdx + 1
        end
    end

    if #tokens < 1 then
        return ""
    end

    for _, v in pairs(tokens) do
        print("TOKEN", v.type, v.val, v.pre, v.len)
    end

    local cmds = {}
    local curCmd = nil
    local invertNextCmd = false
    local token
    local idx = 1
    while idx < #tokens do
        token = tokens[idx]
        if not curCmd then
            curCmd = {
                args = {},
                invert = invertNextCmd,
                stdinMod = nil,
                stdoutMod = nil,
                stderrMod = nil,
                chainToNext = nil,
            }
            invertNextCmd = false
        end

        if token.type == ArgTypeString then
            table.insert(curCmd.args, token.val)
        elseif token.type == ArgTypeOp then
            if token.val == "|" or token.val == "&" or token.val == ";" then
                if token.val ~= ";" then
                    if #curCmd.args < 1 then
                        error("Cannot have " .. token.raw .. " at the start of a command!")
                    end
                    if token.len > 2 then
                        error("Cannot have more than 2 of " .. token.val .. " in a row")
                    end
                    curCmd.chainToNext = token.raw
                end
                -- TODO: Singular | should set stdin up
                table.insert(cmds, curCmd)
                curCmd = nil
            elseif token.val == "!" then
                if #curCmd.args > 0 then
                    error("Cannot have \"" .. token.raw .. "\" in the middle of a command!")
                end
                invertNextCmd = (token.len % 2) == 1
            elseif token.val == "<" or token.val == ">" then
                if #curCmd.args < 1 then
                    error("Cannot redirect stdin/out/err of nothing!")
                end
                idx = idx + 1
                local outFile = tokens[idx]
                if outFile.type ~= ArgTypeString then
                    error("Expected string after " .. token.raw)
                end

                local fileInfo = {
                    name = outFile.val,
                    append = token.len > 1,
                }

                if token.val == "<" then
                    if token.pre and token.pre ~= "" then
                        error("Expected nothing before " .. token.raw)
                    end
                    curCmd.stdinMod = fileInfo
                elseif token.val == ">" then
                    if token.pre == "2" then
                        curCmd.stderrMod = fileInfo
                    elseif token.pre == "1" or token.pre == "" or not token.pre then
                        curCmd.stdoutMod = fileInfo
                    else
                        error("Expected nothing, 1 or 2 before " .. token.raw)
                    end
                end
            end
        end

        idx = idx + 1
    end

    if #curCmd.args > 0 then
        table.insert(cmds, curCmd)
    end

    local function pStdMap(op, v)
        if not v then
            return
        end
        print("REDIR", op, v.name, v.append)
    end
    for _, v in pairs(cmds) do
        print("CMD", v.invert, table.concat(v.args, " "))
        pStdMap("<", v.stdinMod)
        pStdMap(">", v.stdoutMod)
        pStdMap("2>", v.stderrMod)
    end

    --[[
    if cmdHandler.has(tokens[1]) then
        table.insert(tokens, 1, exe)
        table.insert(tokens, 2, "-c")
    end
    local sCmd = cmd.new(tokens)
    sCmd:raiseForBadExit(false)
    -- TODO: Assign multi-command here, assign stdio here
    return function()
        sCmd:run()
    end
    ]]
    return ""
end

return M
