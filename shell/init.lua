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

function shell.parser(cmd, lineNo)
    -- TODO: syntax-aware lua end finding
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

shell.commandSearch = {
    "fox.embed.commands",
    "fox.commands",
    "commands",
}
function shell.runCommand(cmd)
    for _, prefix in ipairs(shell.commandSearch) do
        local status, command = pcall(require, prefix .. "." .. cmd)
        if status then
            return command.run(unpack(shell.args))
        end
    end
    error("No such command: " .. cmd)
end
