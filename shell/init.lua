local config_home = env["XDG_CONFIG_HOME"]
if (not config_home) or config_home == "" then
    config_home = env["HOME"] .. "/.config"
end
local lua_base = config_home .. "/fox/lua"
_G.LUA_BASE = lua_base

package.path = lua_base .. "/modules/?.lua;" .. lua_base .. "/modules/?/init.lua"
package.cpath = ""

local ok, err = pcall(dofile, lua_base .. "/init.lua")
if not ok then
    print("Error loading user init.lua: " .. tostring(err))
end
