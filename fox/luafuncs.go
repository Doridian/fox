package fox

import (
	lua "github.com/yuin/gopher-lua"
)

func luaVersion(L *lua.LState) int {
	L.Push(lua.LString(version))
	return 1
}

func luaGitRev(L *lua.LState) int {
	L.Push(lua.LString(gitrev))
	return 1
}
