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
    pcall(ctx.stdin.close, ctx.stdin)
    pcall(ctx.stdout.close, ctx.stdout)
    pcall(ctx.stderr.close, ctx.stderr)
end

function M.run(ctx, cmd, args)
    local mod, err = getCommand(cmd)
    if err then
        error("Error loading command " .. cmd .. ": " .. err)
    end

    if not mod then
        error("No such command: " .. cmd)
    end

    ctx.name = args[1]
    ctx.stdin = ctx.stdin or pipe.stdin
    ctx.stdout = ctx.stdout or pipe.stdout
    ctx.stderr = ctx.stderr or pipe.stderr
    table.remove(args, 1)
    local exitCode = mod.run(ctx, table.unpack(args))
    closeCtx(ctx)
    return exitCode
end

function M.has(cmd)
    local mod, _ = getCommand(cmd)
    if mod then
        return mod.run or mod.runDirect, mod.runDirect
    end
    return false, false
end

return M
