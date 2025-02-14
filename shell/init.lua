local shell = require("fox.shell")
local Env = require("fox.Env")

local config_home = Env["XDG_CONFIG_HOME"]
if (not config_home) or config_home == "" then
    config_home = Env["HOME"] .. "/.config"
end
local lua_base = config_home .. "/fox/lua"
_G.LUA_BASE = lua_base

package.path = lua_base .. "/modules/?.lua;" .. lua_base .. "/modules/?/init.lua"
package.cpath = ""

local ok, err = pcall(dofile, lua_base .. "/init.lua")
if not ok then
    print("Error loading user init.lua: " .. tostring(err))
end

shell.parsers = shell.parsers or {}
function shell.parsers.lua(cmd)
    cmdLen = cmd:len()
    cmdLastTwo = cmd:sub(cmdLen - 1, cmdLen)
    if cmdLastTwo == "\n\n" then
        return cmd
    end
    return true
end

function shell.parser(cmd)
    if cmd:sub(1, 1) == "!" then
        newLine = cmd:find("\n", 1, true)
        cmdPrefix = cmd:sub(2, newLine - 1)
        if shell.parsers[cmdPrefix] then
            return shell.parsers[cmdPrefix](cmd:sub(newLine + 1))
        end

        print("Unknown parser " .. cmdPrefix)
        return ""
    end

    return false
end
