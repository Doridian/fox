package file

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs"
const LuaTypeName = "File"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"stat":  fileStat,
		"close": fileClose,
		"read":  fileRead,
		"write": fileWrite,
	}))
	mt.RawSetString("__tostring", L.NewFunction(fileToString))
	tbl.RawSetString(LuaTypeName, mt)
}
