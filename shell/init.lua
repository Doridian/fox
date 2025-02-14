local shell = require("fox.shell")
local Env = require("fox.env")
local fs = require("fox.fs")

local configHome = Env["XDG_CONFIG_HOME"]
if (not configHome) or configHome == "" then
    configHome = Env["HOME"] .. "/.config"
end

local baseDir = configHome .. "/fox"
_G.BaseDir = baseDir
fs.mkdirAll(baseDir)

local readlineConfig = shell.getReadlineConfig()
readlineConfig:historyFile(baseDir .. "/history")
shell.readlineConfig(readlineConfig)

package.path = baseDir .. "/modules/?.lua;" .. baseDir .. "/modules/?/init.lua"
package.cpath = ""

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

local ok, err = pcall(dofile, baseDir .. "/init.lua")
if not ok then
    print("Error loading user init.lua: " .. tostring(err))
end
