package vars

import (
	lua "github.com/yuin/gopher-lua"
)

// __index(t, k)
func varsIndex(L *lua.LState) int {
	k := L.CheckString(2)

	v, ok := varTable[k]
	if !ok {
		L.Push(lua.LNil)
		return 1
	}

	L.Push(v)
	return 1
}

// __newindex(t, k, v)
func varsNewIndex(L *lua.LState) int {
	k := L.CheckString(2)
	v := L.CheckString(3)
	varTable[k] = lua.LString(v)

	return 1
}

// __call()
func varsCall(L *lua.LState) int {
	ret := L.NewTable()

	for name, val := range varTable {
		ret.RawSetString(name, val)
	}

	L.Push(ret)
	return 1
}
