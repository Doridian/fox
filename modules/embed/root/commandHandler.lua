local shell = require("go:shell")
local pipe = require("go:pipe")

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

function M.run(cmd, args, direct)
    local mod, err = getCommand(cmd)
    if err then
        error("Error loading command " .. cmd .. ": " .. err)
    end

    if not mod then
        error("No such command: " .. cmd)
    end

    local runFunc
    if direct then
        runFunc = mod.runDirect
    elseif mod.run then
        runFunc = mod.run
    elseif mod.runDirect then
        runFunc = function(...)
            local exitCode, stdout = mod.runDirect(...)
            if stdout then
                pipe.stdout:write(stdout)
            end
            return exitCode
        end
    end
    return runFunc(table.unpack(args))
end

function M.has(cmd)
    local mod, _ = getCommand(cmd)
    if mod then
        return mod.run or mod.runDirect, mod.runDirect
    end
    return false, false
end

return M
