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
	L.ArgError(i, "shellcmd expected")
	return nil, nil
}
