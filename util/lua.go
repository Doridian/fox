package util

import lua "github.com/yuin/gopher-lua"

func RetNil(L *lua.LState) int {
	L.Push(lua.LNil)
	return 1
}
