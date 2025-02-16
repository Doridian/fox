local M =  {}

function M.run(cmdAdd, lineNo, prev)
    local cmd = (prev or "") .. cmdAdd .. "\n"

    if cmdAdd == "" then
        return cmd
    end
    return cmd, true
end

return M
