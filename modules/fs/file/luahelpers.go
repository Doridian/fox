package direntry

import (
	"io/fs"

	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, f fs.File) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = f
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func Push(L *lua.LState, f fs.File) int {
	if f == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(ToUserdata(L, f))
	return 1
}

func Check(L *lua.LState, i int) (fs.File, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(fs.File); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return nil, nil
}
