local DEFAULT_MODS = {
    "env",
    "cmd",
    "pipe",
}
for _, m in pairs(DEFAULT_MODS) do
    _G[m] = require(m)
end

local embedded = require("embedded")
table.insert(package.loaders, embedded.loader)
