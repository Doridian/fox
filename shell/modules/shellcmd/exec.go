package shellcmd

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func handleCmdExit(L *lua.LState, exitCode int, c *Cmd, ud *lua.LUserData, err error) int {
	c.releaseStdio()

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

func doWaitCmd(L *lua.LState, c *Cmd, ud *lua.LUserData) int {
	if c == nil {
		return 0
	}
	pipeErr := c.waitStdio()
	err := c.gocmd.Wait()
	if err == nil {
		err = pipeErr
	}
	return handleCmdExit(L, c.gocmd.ProcessState.ExitCode(), c, ud, err)
}

func doWait(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	return doWaitCmd(L, c, ud)
}

func doRun(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	err := c.prepareAndStart()
	if err != nil {
		return handleCmdExit(L, 1, c, ud, err)
	}
	return doWaitCmd(L, c, ud)
}

func doStart(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	err := c.prepareAndStart()
	if err != nil {
		return handleCmdExit(L, 1, c, ud, err)
	}
	L.Push(ud)
	return 1
}

func (c *Cmd) prepareAndStart() error {
	if err := c.setupStdio(); err != nil {
		return err
	}

	return c.gocmd.Start()
}

func (c *Cmd) prepareAndRun() error {
	if err := c.prepareAndStart(); err != nil {
		return err
	}

	pipeErr := c.waitStdio()
	err := c.gocmd.Wait()
	if err == nil {
		err = pipeErr
	}
	return err
}
