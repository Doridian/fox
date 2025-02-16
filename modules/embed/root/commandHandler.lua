local shell = require("go:shell")

local M = {}
local cmdCache = {}

M.search = {
    "commands",
    "go:commands",
    "embed:commands",
}
local function getCommand(cmd)
    if cmdCache[cmd] then
        return cmdCache[cmd]
    end

    for _, prefix in pairs(M.search) do
        local ok, mod = pcall(require, prefix .. "." .. cmd)
        if ok then
            cmdCache[cmd] = mod
            return mod
        end
    end
    return nil
end

function M.run(cmd, args)
    local mod = getCommand(cmd)
    if mod then
        return mod.run(table.unpack(args))
    end
    error("No such command: " .. cmd)
end

function M.has(cmd)
    if getCommand(cmd) then
        return true
    end
    return false
end

return M
