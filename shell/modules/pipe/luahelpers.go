package pipe

import lua "github.com/yuin/gopher-lua"

func pushPipe(L *lua.LState, pipe *Pipe) int {
	if pipe == nil {
		L.Push(lua.LNil)
		return 1
	}

	ud := L.NewUserData()
	ud.Value = pipe
	L.SetMetatable(ud, L.GetTypeMetatable(luaShellPipeType))
	L.Push(ud)
	return 1
}

func CheckPipe[K interface{}](L *lua.LState, i int, allowNil bool) (bool, *Pipe, *lua.LUserData) {
	if L.Get(i) == lua.LNil && allowNil {
		return true, nil, nil
	}

	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*Pipe); ok {
		return true, v, ud
	}

	if allowNil {
		L.ArgError(i, "pipe or nil expected")
	} else {
		L.ArgError(i, "pipe expected")
	}

	return false, nil, nil
}
