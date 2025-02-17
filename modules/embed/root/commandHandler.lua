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

function M.closeCtx(ctx)
    for _, v in pairs(ctx) do
        if v and v.close then
            pcall(v.close, v)
        end
    end
end

function M.run(ctx, cmd, args)
    local mod, err = getCommand(cmd)
    if err then
        error("Error loading command " .. cmd .. ": " .. err)
    end

    if not mod then
        error("No such command: " .. cmd)
    end

    local exitCode = mod.run(ctx, table.unpack(args))
    M.closeCtx(ctx)
    return exitCode
end

function M.get(cmd)
    return getCommand(cmd)
end

return M
