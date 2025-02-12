package pipe

import lua "github.com/yuin/gopher-lua"

func Make(L *lua.LState, pipe *Pipe) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = pipe
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func Push(L *lua.LState, pipe *Pipe) int {
	if pipe == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(Make(L, pipe))
	return 1
}

func Check(L *lua.LState, i int, allowNil bool) (bool, *Pipe, *lua.LUserData) {
	if L.Get(i) == lua.LNil && allowNil {
		return true, nil, nil
	}

	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*Pipe); ok {
		return true, v, ud
	}

	if allowNil {
		L.ArgError(i, LuaType+" or nil expected")
	} else {
		L.ArgError(i, LuaType+" expected")
	}

	return false, nil, nil
}
