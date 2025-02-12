package fileinfo

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs.fileinfo"
const LuaType = LuaName + ":FileInfo"

func Load(L *lua.LState) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"name":    fiName,
		"size":    fiSize,
		"mode":    fiMode,
		"modTime": fiModTime,
		"isDir":   fiIsDir,
	}))
}
