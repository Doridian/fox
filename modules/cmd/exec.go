package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func handleCmdExitNoLock(L *lua.LState, exitCode int, c *Cmd, ud *lua.LUserData, err error) int {
	c.releaseStdioNoLock()

	if err != nil && exitCode == 0 {
		exitCode = 1
	}
	exitCodeL := lua.LNumber(exitCode)

	L.SetGlobal("_LAST_EXIT_CODE", exitCodeL)
	if c.ErrorPropagation {
		if err == nil {
			err = fmt.Errorf("command exited with code %d", exitCode)
		}
		L.RaiseError("%s", err.Error())
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

func doWaitCmdNoLock(L *lua.LState, c *Cmd, ud *lua.LUserData) int {
	pipeErr := c.waitStdio()
	err := c.gocmd.Wait()
	if err == nil {
		err = pipeErr
	}
	return handleCmdExitNoLock(L, c.gocmd.ProcessState.ExitCode(), c, ud, err)
}

func doWait(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	return doWaitCmdNoLock(L, c, ud)
}

func doRun(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	err := c.prepareAndStartNoLock()
	if err != nil {
		return handleCmdExitNoLock(L, 1, c, ud, err)
	}
	return doWaitCmdNoLock(L, c, ud)
}

func doStart(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	defer c.lock.RUnlock()

	err := c.prepareAndStartNoLock()
	if err != nil {
		return handleCmdExitNoLock(L, 1, c, ud, err)
	}
	L.Push(ud)
	return 1
}

func (c *Cmd) prepareAndStartNoLock() error {
	var err error

	path := c.gocmd.Args[0]
	if c.AutoLookPath && !strings.ContainsRune(path, '/') {
		path, err = exec.LookPath(path)
		if err != nil {
			return err
		}
	}
	c.gocmd.Path = path

	err = c.setupStdio()
	if err != nil {
		return err
	}

	return c.gocmd.Start()
}

func (c *Cmd) ensureRan() error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.gocmd.Process == nil {
		if err := c.prepareAndStartNoLock(); err != nil {
			return err
		}
	}

	pipeErr := c.waitStdio()
	err := c.gocmd.Wait()
	if err == nil {
		err = pipeErr
	}
	return err
}
