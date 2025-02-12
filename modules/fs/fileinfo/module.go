package fileinfo

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs"
const LuaTypeName = "FileInfo"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"name":    fiName,
		"size":    fiSize,
		"mode":    fiMode,
		"modTime": fiModTime,
		"isDir":   fiIsDir,
	}))
	mt.RawSetString("__tostring", L.NewFunction(fiToString))
	tbl.RawSetString(LuaTypeName, mt)
}
