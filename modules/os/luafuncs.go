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

func osChdir(L *lua.LState) int {
	dir := L.ToString(1)
	err := goos.Chdir(dir)
	if err != nil {
		L.RaiseError("os.Chdir(%q) failed: %s", dir, err)
	}
	return 0
}

func osGetwd(L *lua.LState) int {
	pwd, err := goos.Getwd()
	if err != nil {
		L.RaiseError("os.Getwd() failed: %s", err)
		return 0
	}
	L.Push(lua.LString(pwd))
	return 1
}
