package direntry

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:fox.fs"
const LuaTypeName = "DirEntry"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"name":  deName,
		"isDir": deIsDir,
		"type":  deType,
		"info":  deInfo,
	}))
	mt.RawSetString("__tostring", L.NewFunction(deToString))
	tbl.RawSetString(LuaTypeName, mt)
}
