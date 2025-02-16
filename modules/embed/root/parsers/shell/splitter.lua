local tokenizer = require("embed:parsers.shell.tokenizer")

local M = {}

M.RedirTypeFile = 1
M.RedirTypeCmd = 2

function M.run(tokens)
    local cmds = {}
    local curCmd = nil
    local invertNextCmd = false
    local stdinNextCmd = nil
    local token
    local idx = 1
    while idx <= #tokens do
        token = tokens[idx]
        print("TOKEN", token.type, token.val, token.pre, token.len)

        if not curCmd then
            curCmd = {
                args = {},
                invert = invertNextCmd,
                stdin = stdinNextCmd,
                stdout = nil,
                stderr = nil,
                chainToNext = nil,
                background = false,
            }
            stdinNextCmd = nil
            invertNextCmd = false
        end

        if token.type == tokenizer.ArgTypeString then
            table.insert(curCmd.args, token.val)
        elseif token.type == tokenizer.ArgTypeOp then
            if token.val == "|" or token.val == "&" or token.val == ";" then
                if token.val == ";" then
                    if #curCmd.args < 1 and (curCmd.invert or curCmd.stdin) then
                        error("Cannot have " .. token.val .. " after " .. curCmd.invert and "!" or "|")
                    end
                elseif #curCmd.args < 1 then
                    error("Cannot have " .. token.raw .. " at the start of a command!")
                else
                    if token.val == "|" and token.len == 1 then
                        stdinNextCmd = {
                            type = M.RedirTypeCmd,
                            cmd = curCmd,
                        }
                        if curCmd.background then
                            error("Cannot pipe (" .. token.raw .. ") after background command (&)")
                        end
                    elseif token.val == "&" and token.len == 1 then
                        curCmd.background = true
                    else
                        if token.len > 2 then
                            error("Cannot have more than 2 of " .. token.val .. " in a row")
                        end
                        curCmd.chainToNext = token.raw
                    end
                end
                if #curCmd.args > 0 then
                    table.insert(cmds, curCmd)
                end
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
                if outFile.type ~= tokenizer.ArgTypeString then
                    error("Expected string after " .. token.raw)
                end

                local fileInfo = {
                    type = M.RedirTypeFile,
                    name = outFile.val,
                    append = token.len > 1,
                }

                if token.val == "<" then
                    if token.pre and token.pre ~= "" then
                        error("Expected nothing before " .. token.raw)
                    end
                    curCmd.stdin = fileInfo
                elseif token.val == ">" then
                    if token.pre == "2" then
                        curCmd.stderr = fileInfo
                    elseif token.pre == "1" or token.pre == "" or not token.pre then
                        curCmd.stdout = fileInfo
                    else
                        error("Expected nothing, 1 or 2 before " .. token.raw)
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
