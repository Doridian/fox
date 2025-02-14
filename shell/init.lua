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

package.path = baseDir .. "/modules/?.lua;" .. baseDir .. "/modules/?/init.lua"
package.cpath = ""

shell.parsers = {}
function shell.parsers.lua(cmd, lineNo)
    local cmdLen = cmd:len()
    if cmd:sub(cmdLen - 1, cmdLen) == "\n\n" then
        return cmd
    end
    return true
end

shell.commands = {}
function shell.parsers.cmd(cmd, lineNo)
    local parsed, promptOverride = shell.defaultShellParser(cmd, lineNo)
    if (not parsed) or parsed == true or parsed == "" then
        return parsed, promptOverride
    end

    -- TODO: Parse CLI-like language
    return parsed
end
shell.parsers.default = shell.parsers.cmd

function shell.parser(cmd, lineNo)
    if cmd:sub(1, 1) == "!" then
        local newLine = cmd:find("\n", 1, true)
        local cmdPrefix = cmd:sub(2, newLine - 1)
        if shell.parsers[cmdPrefix] then
            return shell.parsers[cmdPrefix](cmd:sub(newLine + 1), lineNo)
        end

        print("Unknown parser " .. cmdPrefix)
        return ""
    end

    local defParser = shell.parsers.default
    if defParser then
        return defParser(cmd, lineNo)
    end
    return false
end

function shell.setHistoryFile(file)
    local readlineConfig = shell.getReadlineConfig()
    readlineConfig:historyFile(file)
    shell.readlineConfig(readlineConfig)
end
shell.setHistoryFile(baseDir .. "/history")

local initLua = baseDir .. "/init.lua"
if fs.stat(initLua) then
    local ok, err = pcall(dofile, initLua)
    if not ok then
        print("Error loading user init.lua: " .. tostring(err))
    end
end
