local tokenizer = require("embed:parsers.shell.tokenizer")

local M = {}

M.RedirTypeFile = 1
M.RedirTypeCmd = 2
M.RedirTypeRefer = 3

function M.run(tokens)
    local cmds = {}
    local curCmd = nil
    local stdinNextCmd = nil
    local token
    local idx = 1
    while idx <= #tokens do
        token = tokens[idx]

        if not curCmd then
            curCmd = {
                args = {},
                invert = false,
                stdin = stdinNextCmd,
                stdout = nil,
                stderr = nil,
                chainToNext = nil,
                background = false,
            }
            stdinNextCmd = nil
        end

        if token.type == tokenizer.ArgTypeString or token.type == tokenizer.ArgTypeStringFunc then
            table.insert(curCmd.args, token)
        elseif token.type == tokenizer.ArgTypeOp then
            if token.value == "|" or token.value == "&" or token.value == ";" then
                if token.value == ";" then
                    if #curCmd.args < 1 and (curCmd.invert or curCmd.stdin) then
                        return nil, "Cannot have " .. token.value .. " after " .. curCmd.invert and "!" or "|"
                    end
                elseif #curCmd.args < 1 then
                    return nil, "cannot have " .. token.raw .. " at the start of a command!"
                else
                    if token.value == "|" and token.len == 1 then
                        stdinNextCmd = {
                            type = M.RedirTypeCmd,
                            cmd = curCmd,
                        }
                        if curCmd.background then
                            return nil, "cannot pipe (" .. token.raw .. ") after background command (&)"
                        end
                    elseif token.value == "&" and token.len == 1 then
                        curCmd.background = true
                    else
                        if token.len > 2 then
                            return nil, "cannot have more than 2 of " .. token.value .. " in a row"
                        end
                        -- Must have || or && here
                        if curCmd.background then
                            return nil, "cannot chain (" .. token.raw .. ") to background command (&)"
                        end
                        curCmd.chainToNext = token.raw
                    end
                end
                if #curCmd.args > 0 then
                    table.insert(cmds, curCmd)
                end
                curCmd = nil
            elseif token.value == "!" then
                if #curCmd.args > 0 then
                    return nil, "cannot have \"" .. token.raw .. "\" in the middle of a command"
                end
                local invLen = token.len
                if curCmd.invert then
                    invLen = invLen + 1
                end
                curCmd.invert = (invLen % 2) == 1
            elseif token.value == "<" or token.value == ">" then
                if #curCmd.args < 1 then
                    return nil, "cannot redirect stdin/out/err of nothing"
                end
                if token.value == ">" and token.len > 2 then
                    return nil, "cannot have more than 2 of " .. token.value .. " in a row"
                elseif token.value == "<" and token.len > 1 then
                    return nil, "cannot have more than 1 of " .. token.value .. " in a row"
                end

                idx = idx + 1
                local outFile = tokens[idx]
                if outFile.type ~= tokenizer.ArgTypeString and outFile.type ~= tokenizer.ArgTypeStringFunc then
                    return nil, "expected string after " .. token.raw
                end

                local fileInfo = {
                    type = M.RedirTypeFile,
                    name = outFile,
                    append = token.len > 1,
                }

                if token.hasAmpersand then
                    local referTo = tonumber(outFile.val)
                    if referTo == 1 then
                        fileInfo.type = M.RedirTypeRefer
                        fileInfo.ref = "stdout"
                    elseif referTo == 2 then
                        fileInfo.type = M.RedirTypeRefer
                        fileInfo.ref = "stderr"
                    else
                        return nil, "&ref to invalid stream: " .. outFile.val
                    end
                end

                if token.value == "<" then
                    if token.pre and token.pre ~= "" then
                        return nil, "expected nothing before " .. token.raw
                    end
                    curCmd.stdin = fileInfo
                elseif token.value == ">" then
                    if token.pre == "2" then
                        curCmd.stderr = fileInfo
                    elseif token.pre == "1" or token.pre == "" or not token.pre then
                        curCmd.stdout = fileInfo
                    else
                        return nil, "expected nothing, 1 or 2 before " .. token.raw
                    end
                end
            end
        end

        idx = idx + 1
    end

    if curCmd and #curCmd.args > 0 then
        table.insert(cmds, curCmd)
    end

    return cmds
end

return M
