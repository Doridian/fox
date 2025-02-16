package pipe

import (
	lua "github.com/yuin/gopher-lua"
)

func pipeToString(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}

	L.Push(lua.LString(p.ToString()))
	return 1
}

func pipeIsNull(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	L.Push(lua.LBool(p.IsNull()))
	return 1
}
