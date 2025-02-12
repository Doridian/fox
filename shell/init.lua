local GLOBAL_MODS = {
    "env",
    "cmd",
    "pipe",
    "embedded",
}
for _, m in pairs(GLOBAL_MODS) do
    _G[m] = require("fox." .. m)
end

table.insert(package.loaders, embedded.prefixLoader)
package.cpath = ""

-- TODO: Respect XDG_CONFIG_HOME
local lua_base = env["HOME"] .. "/.config/fox/lua"
_G.LUA_BASE = lua_base
package.path = lua_base .. "/modules/?.lua;" .. lua_base .. "/modules/?/init.lua"

local ok, err = pcall(dofile, lua_base .. "/init.lua")
if not ok then
    print("Error loading user init.lua: " .. err)
end
