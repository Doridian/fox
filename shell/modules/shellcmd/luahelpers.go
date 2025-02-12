package shellcmd

import lua "github.com/yuin/gopher-lua"

func pushShellCmd(L *lua.LState, cmd *ShellCmd) int {
	if cmd == nil {
		L.Push(lua.LNil)
		return 1
	}

	ud := L.NewUserData()
	ud.Value = cmd
	L.SetMetatable(ud, L.GetTypeMetatable(luaShellCmdType))
	L.Push(ud)
	return 1
}

func checkShellCmd(L *lua.LState, i int) (*ShellCmd, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*ShellCmd); ok {
		return v, ud
	}
	L.ArgError(i, "ShellCmd expected")
	return nil, nil
}

func pushShellPipe(L *lua.LState, pipe *Pipe) int {
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

func checkPipe(L *lua.LState, i int, allowNil bool) (bool, *Pipe, *lua.LUserData) {
	if L.Get(i) == lua.LNil && allowNil {
		return true, nil, nil
	}

	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*Pipe); ok {
		return true, v, ud
	}

	if allowNil {
		L.ArgError(i, "Pipe or nil expected")
	} else {
		L.ArgError(i, "Pipe expected")
	}

	return false, nil, nil
}
