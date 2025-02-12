package direntry

import (
	"io/fs"

	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, de fs.DirEntry) *lua.LUserData {
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
	L.Push(ToUserdata(L, de))
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

func ArrayToUserdata(L *lua.LState, dirents []fs.DirEntry) *lua.LTable {
	ret := L.NewTable()
	for _, de := range dirents {
		ret.Append(ToUserdata(L, de))
	}
	return ret
}

func PushArray(L *lua.LState, dirents []fs.DirEntry) int {
	L.Push(ArrayToUserdata(L, dirents))
	return 1
}
