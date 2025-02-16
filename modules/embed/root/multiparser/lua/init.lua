local M =  {}

function M.run(cmdAdd, lineNo, prev)
    local cmd = prev .. cmdAdd

    if cmd:sub(#cmd - 1, #cmd) == "\n\n" then
        return cmd
    end
    return cmd, true
end

return M
