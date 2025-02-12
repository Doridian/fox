package pipe

import (
	lua "github.com/yuin/gopher-lua"
)

func pipeCanWrite(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	L.Push(lua.LBool(p.CanWrite()))
	return 1
}

func pipeCanRead(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	L.Push(lua.LBool(p.CanRead()))
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

func pipeToString(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}

	L.Push(lua.LString(p.ToString()))
	return 1
}
