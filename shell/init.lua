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
function shell.runCommand(cmd)
    for _, prefix in pairs(shell.commandSearch) do
        local ok, mod = pcall(require, prefix .. "." .. cmd)
        if ok then
            return mod.run(table.unpack(shell.args))
        end
    end
    error("No such command: " .. cmd)
end
function shell.hasCommand(cmd)
    for _, prefix in pairs(shell.commandSearch) do
        local ok, _ = pcall(require, prefix .. "." .. cmd)
        if ok then
            return true
        end
    end
    return false
end

local mp = require("embed:multiparser")
shell.parser = mp.run
mp.parsers.default = mp.parsers.shell

local initLua = baseDir .. "/init.lua"
if fs.stat(initLua) then
    local ok, err = pcall(dofile, initLua)
    if not ok then
        print("Error loading user init.lua: " .. tostring(err))
    end
end
