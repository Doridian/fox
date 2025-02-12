-- TODO: Respect XDG_CONFIG_HOME
local lua_base = tostring(env["HOME"]) .. "/.config/fox/lua"
_G.LUA_BASE = lua_base
package.path = lua_base .. "/modules/?.lua;" .. lua_base .. "/modules/?/init.lua"
package.cpath = ""

local ok, err = pcall(dofile, lua_base .. "/init.lua")
if not ok then
    print("Error loading user init.lua: " .. tostring(err))
end
