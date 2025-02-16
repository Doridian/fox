local shell = require("go:shell")
local env = require("go:env")
local fs = require("go:fs")

---@diagnostic disable-next-line: deprecated
table.unpack = table.unpack or _G.unpack

local configHome = env["XDG_CONFIG_HOME"]
if (not configHome) or configHome == "" then
    configHome = env["HOME"] .. "/.config"
end

local baseDir = configHome .. "/fox"
_G.BaseDir = baseDir
fs.mkdirAll(baseDir)

package.path = baseDir .. "/modules/?.lua;" .. baseDir .. "/modules/?/init.lua"
package.cpath = ""

function shell.setHistoryFile(file)
    local readlineConfig = shell.getReadlineConfig()
    readlineConfig:historyFile(file)
    shell.readlineConfig(readlineConfig)
end
shell.setHistoryFile(baseDir .. "/history")

shell.commandSearch = {
    "commands",
    "go:commands",
    "embed:commands",
}
local function getCommand(cmd)
    for _, prefix in pairs(shell.commandSearch) do
        local ok, mod = pcall(require, prefix .. "." .. cmd)
        if ok then
            return mod
        end
    end
    return nil
end
function shell.runCommand(cmd)
    local mod = getCommand(cmd)
    if mod then
        return mod.run(table.unpack(shell.args()))
    end
    error("No such command: " .. cmd)
end
function shell.hasCommand(cmd)
    if getCommand(cmd) then
        return true
    end
    return false
end

if shell.interactive() then
    local mp = require("embed:multiparser")
    mp.defaultParser = "shell"
    shell.parser = mp.run
end

local initLua = baseDir .. "/init.lua"
if fs.stat(initLua) then
    local ok, err = pcall(dofile, initLua)
    if not ok then
        print("Error loading user init.lua: " .. tostring(err))
    end
end
