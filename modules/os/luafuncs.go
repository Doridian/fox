package os

import (
	goos "os"

	lua "github.com/yuin/gopher-lua"
)

func osExecutable(L *lua.LState) int {
	exe, err := goos.Executable()
	if err != nil {
		L.RaiseError("os.Executable() failed: %s", err)
		return 0
	}
	L.Push(lua.LString(exe))
	return 1
}
