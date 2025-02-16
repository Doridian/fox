local os = require("go:os")
local tokenizer = require("embed:parsers.shell.tokenizer")
local splitter = require("embed:parsers.shell.splitter")

local exe = os.executable()

--[[
    TODO:
    - Need to do command splitting at parse time before variable/glob interp
        - | || & && ;
    - Need todo stdio redirections
]]

local M = {}

function M.run(strAdd, lineNo, prev)
    local parsed = (prev or "") .. strAdd .. "\n"

    if strAdd:sub(#strAdd, #strAdd) == "\\" then
        return parsed:sub(1, #parsed - 2) .. "\n", true
    end

    local tokens = tokenizer.run(parsed)
    local cmds = splitter.run(tokens)

    local function pStdMap(op, v)
        if not v then
            return
        end
        print("REDIR", op, v.type, v.name, v.cmd, v.append)
    end
    for _, v in pairs(cmds) do
        print("CMD", v.invert, v.background, table.concat(v.args, " "))
        pStdMap("STDIN", v.stdin)
        pStdMap("STDOUT", v.stdout)
        pStdMap("STDERR", v.stderr)
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
