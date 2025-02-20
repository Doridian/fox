local shell = require("go:shell")
local env = require("go:env")
local fs = require("go:fs")

local p = print
function _G.print(...)
    shell.stdout:print(...)
end

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

local cmdHandler = require("embed:commandHandler")
shell.runCommand = cmdHandler.run

if shell.interactive() then
    function shell.setHistoryFile(file)
        local readlineConfig = shell.getReadlineConfig()
        readlineConfig:historyFile(file)
        shell.readlineConfig(readlineConfig)
    end
    shell.setHistoryFile(baseDir .. "/history")

    local mp = require("embed:multiParser")
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
