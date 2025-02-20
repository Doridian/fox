local M = {}

local commandCache = {}

M.search = {
    "commands",
    "embed:commands",
}
local function getCommand(name)
    local c = commandCache[name]
    if c then
        if c.err then
            return nil, c.err
        end
        return c
    end

    local errs = {}
    for _, prefix in pairs(M.search) do
        local modName = prefix .. "." .. name
        local ok, mod = pcall(require, modName)
        if ok then
            commandCache[name] = mod
            return mod
        elseif not mod then
            table.insert(errs, "require() did not return table for ".. modName)
        else
            table.insert(errs, mod)
        end
    end

    local err = table.concat(errs, "\n")
    commandCache[name] = {
        err = err,
    }
    return nil, err
end

local function mustGetCommand(name)
    local p, err = getCommand(name)
    if not p then
        error("Error loading command " .. name .. ": " .. err)
    end
    return p
end

function M.get(name)
    return mustGetCommand(name)
end

function M.run(name, ...)
    local cmd = mustGetCommand(name)
    return cmd.run(name, ...)
end

return M
