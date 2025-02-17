local M = {}
local cmdCache = {}

M.search = {
    "commands",
    "embed:commands",
}
function M.get(cmd)
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
    if not ctx then
        return
    end

    for _, v in pairs(ctx) do
        if v and v.close then
            pcall(v.close, v)
        end
    end
end

function M.run(ctx, cmd, args)
    local mod, err = M.get(cmd)
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

return M
