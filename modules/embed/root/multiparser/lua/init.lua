local M =  {}

function M.run(cmd, lineNo)
    if cmd:sub(#cmd - 1, #cmd) == "\n\n" then
        return cmd
    end
    return true
end

return M
