package direntry

import (
	"io/fs"

	lua "github.com/yuin/gopher-lua"
)

func Make(L *lua.LState, de fs.DirEntry) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = de
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func Push(L *lua.LState, de fs.DirEntry) int {
	if de == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(Make(L, de))
	return 1
}

func Check(L *lua.LState, i int) (fs.DirEntry, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(fs.DirEntry); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return nil, nil
}
