local GLOBAL_MODS = {
    "env",
    "cmd",
    "pipe",
}
for _, m in pairs(GLOBAL_MODS) do
    _G[m] = require("fox." .. m)
end

local embedded = require("fox.embedded")
print(require("fox.embedded"), require("fox.embedded"))
table.insert(package.loaders, embedded.loader)

print(pcall(require, "abcd.test"))
