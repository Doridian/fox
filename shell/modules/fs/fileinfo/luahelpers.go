package fileinfo

import (
	"io/fs"

	lua "github.com/yuin/gopher-lua"
)

func Make(L *lua.LState, fi fs.FileInfo) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = fi
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func Push(L *lua.LState, fi fs.FileInfo) int {
	if fi == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(Make(L, fi))
	return 1
}

func Check(L *lua.LState, i int) (fs.FileInfo, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(fs.FileInfo); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return nil, nil
}
