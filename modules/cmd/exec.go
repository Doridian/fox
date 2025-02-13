package cmd

import (
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func handleCmdExitNoLock(L *lua.LState, exitCode int, c *Cmd, ud *lua.LUserData) int {
	c.releaseStdioNoLock()

	exitCodeL := lua.LNumber(exitCode)

	L.SetGlobal("_LAST_EXIT_CODE", exitCodeL)
	if c.ErrorPropagation && exitCode != 0 {
		L.RaiseError("command exited with code %d", exitCode)
		return 0
	}

	L.Push(ud)
	L.Push(exitCodeL)
	return 2
}

func (c *Cmd) doWaitCmdNoLock() {
	c.awaited = true
	c.waitSync.Wait()
}

func doWaitCmdNoLock(L *lua.LState, c *Cmd, ud *lua.LUserData) int {
	c.doWaitCmdNoLock()
	return handleCmdExitNoLock(L, c.gocmd.ProcessState.ExitCode(), c, ud)
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
	return c.doRun(L, ud)
}

func (c *Cmd) doRun(L *lua.LState, ud *lua.LUserData) int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	err := c.prepareAndStartNoLock(true)
	if err != nil {
		return handleCmdExitNoLock(L, -10002, c, ud)
	}
	return doWaitCmdNoLock(L, c, ud)
}

func doStart(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}
	return c.doStart(L, ud)
}

func (c *Cmd) doStart(L *lua.LState, ud *lua.LUserData) int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	err := c.prepareAndStartNoLock(false)
	if err != nil {
		return handleCmdExitNoLock(L, -10002, c, ud)
	}
	L.Push(ud)
	return 1
}

func (c *Cmd) prepareAndStartNoLock(defaultStdin bool) error {
	var err error

	path := c.gocmd.Args[0]
	if c.AutoLookPath && !strings.ContainsRune(path, '/') {
		path, err = exec.LookPath(path)
		if err != nil {
			return err
		}
	}
	c.gocmd.Path = path

	err = c.setupStdio(defaultStdin)
	if err != nil {
		return err
	}

	err = c.gocmd.Start()
	if err != nil {
		return err
	}

	c.waitSync.Add(1)
	c.mod.addCmd(c)
	go func() {
		c.waitStdio()
		c.gocmd.Wait()
		c.waitSync.Done()
	}()

	return nil
}

func (c *Cmd) ensureRan() error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.gocmd.Process == nil {
		if err := c.prepareAndStartNoLock(true); err != nil {
			return err
		}
	}

	c.doWaitCmdNoLock()
	return nil
}

func doKill(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	if c.gocmd.Process != nil {
		c.gocmd.Process.Kill()
	}

	L.Push(ud)
	return 1
}
