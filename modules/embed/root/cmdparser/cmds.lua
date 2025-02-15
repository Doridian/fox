local M = {}

M.map = {}

function M.loadDirect(name, func)
    M.map[name] = func
end

function M.load(mod)
    M.map[mod.name] = mod.run
end

function M.require(name)
    M.load(require(name))
end

function M.get(name)
    return M.map[name]
end

local function requireBuiltIn(name)
    M.require("fox.embed.cmdparser.commands." .. name)
end
requireBuiltIn("exit")

return M
