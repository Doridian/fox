package shellcmd

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func handleCmdError(L *lua.LState, exitCode int, c *ShellCmd, ud *lua.LUserData, err error) int {
	if err != nil && exitCode == 0 {
		exitCode = 1
	}
	exitCodeL := lua.LNumber(exitCode)

	L.SetGlobal("_LAST_EXIT_CODE", exitCodeL)
	if c.ErrorPropagation {
		if err == nil {
			err = fmt.Errorf("command exited with code %d", exitCode)
		}
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	L.Push(ud)
	L.Push(exitCodeL)
	if err != nil {
		L.Push(lua.LString(err.Error()))
	} else {
		L.Push(lua.LNil)
	}
	return 3
}

func doWaitCmd(L *lua.LState, c *ShellCmd, ud *lua.LUserData) int {
	if c == nil {
		return 0
	}
	err := c.Gocmd.Wait()
	return handleCmdError(L, c.Gocmd.ProcessState.ExitCode(), c, ud, err)
}

func doWait(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	return doWaitCmd(L, c, ud)
}

func doRun(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	err := c.Gocmd.Start()
	if err != nil {
		return handleCmdError(L, 1, c, ud, err)
	}
	return doWaitCmd(L, c, ud)
}

func doStart(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	err := c.Gocmd.Start()
	if err != nil {
		return handleCmdError(L, 1, c, ud, err)
	}
	L.Push(ud)
	return 1
}
