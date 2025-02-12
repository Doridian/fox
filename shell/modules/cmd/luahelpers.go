package cmd

import lua "github.com/yuin/gopher-lua"

func pushCmd(L *lua.LState, cmd *Cmd) int {
	if cmd == nil {
		L.Push(lua.LNil)
		return 1
	}

	ud := L.NewUserData()
	ud.Value = cmd
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	L.Push(ud)
	return 1
}

func checkCmd(L *lua.LState, i int) (*Cmd, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*Cmd); ok {
		return v, ud
	}
	L.ArgError(i, LuaType+" expected")
	return nil, nil
}
