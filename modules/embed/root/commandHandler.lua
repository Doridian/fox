local shell = require("go:shell")

local M = {}
local cmdCache = {}

M.search = {
    "commands",
    "embed:commands",
}
local function getCommand(cmd)
    local c = cmdCache[cmd]
    if c then
        if c.err then
            return nil, c.err
        end
        return c
    end

    local errs = {}
    for _, prefix in pairs(M.search) do
        local ok, mod = pcall(require, prefix .. "." .. cmd)
        if ok then
            cmdCache[cmd] = mod
            return mod
        else
            table.insert(errs, mod)
        end
    end

    local err = table.concat(errs, "\n")
    cmdCache[cmd] = {
        err = err,
    }
    return nil, err
end

function M.run(cmd, args)
    local mod, err = getCommand(cmd)
    if err then
        error("Error loading command " .. cmd .. ": " .. err)
    end

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
