package info

import (
	lua "github.com/yuin/gopher-lua"
)

// __index(t, k)
func infoIndex(L *lua.LState) int {
	k := L.CheckString(2)

	v, ok := infoTable[k]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(v)
	return 1
}

// __call()
func infoCall(L *lua.LState) int {
	ret := L.NewTable()

	for name, val := range infoTable {
		ret.RawSetString(name, val)
	}

	L.Push(ret)
	return 1
}
