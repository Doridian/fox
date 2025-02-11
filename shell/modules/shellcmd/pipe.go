package shellcmd

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

func getSetStdout(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		c2, _ := checkShellCmd(L, 2)
		if c2 == nil {
			return 0
		}
		c.Stdout = c2
		L.Push(ud)
		return 1
	}

	return pushShellCmd(L, c.Stdout)
}

func getSetStderr(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		c2, _ := checkShellCmd(L, 2)
		if c2 == nil {
			return 0
		}
		c.Stderr = c2
		L.Push(ud)
		return 1
	}

	return pushShellCmd(L, c.Stderr)
}

func getSetStdin(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		c2, _ := checkShellCmd(L, 2)
		if c2 == nil {
			return 0
		}
		c.Stdin = c2
		L.Push(ud)
		return 1
	}

	return pushShellCmd(L, c.Stdin)
}

func (c *ShellCmd) prepareAndStart() error {
	c.Gocmd.Stderr = os.Stderr
	c.Gocmd.Stdout = os.Stdout
	c.Gocmd.Stdin = os.Stdin

	return c.Gocmd.Start()
}
