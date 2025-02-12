package direntry

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.fs.direntry"
const LuaType = LuaName + ":DirEntry"

func Load(L *lua.LState) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"name":  deName,
		"isDir": deIsDir,
		"type":  deType,
		"info":  deInfo,
	}))
}
