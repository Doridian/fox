package cmd

import (
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func handleCmdExitNoLock(L *lua.LState, nonExitError error, exitCode int, c *Cmd) int {
	_ = c.releaseStdioNoLock()

	exitCodeL := lua.LNumber(exitCode)

	if c.RaiseForBadExit && exitCode != 0 {
		L.RaiseError("command exited with code %d", exitCode)
		return 0
	}

	if nonExitError != nil {
		L.RaiseError("%v", nonExitError)
		return 0
	}

	L.Push(exitCodeL)
	return 1
}

func (c *Cmd) doWaitCmdNoLock() {
	c.awaited = true
	c.waitSync.Wait()
}

func doWaitCmdNoLock(L *lua.LState, c *Cmd) int {
	c.doWaitCmdNoLock()
	return handleCmdExitNoLock(L, nil, c.gocmd.ProcessState.ExitCode(), c)
}

func doWait(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	defer c.lock.RUnlock()
	return doWaitCmdNoLock(L, c)
}

func doRun(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}
	return c.doRun(L)
}

func (c *Cmd) doRun(L *lua.LState) int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	err := c.prepareAndStartNoLock(L, true)
	if err != nil {
		return handleCmdExitNoLock(L, err, ExitCodeProcessCouldNotStart, c)
	}
	return doWaitCmdNoLock(L, c)
}

func doStart(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}
	return c.doStart(L)
}

func (c *Cmd) doStart(L *lua.LState) int {
	c.lock.RLock()
	defer c.lock.RUnlock()

	err := c.prepareAndStartNoLock(L, false)
	if err != nil {
		return handleCmdExitNoLock(L, err, ExitCodeProcessCouldNotStart, c)
	}
	return 0
}

func (c *Cmd) prepareAndStartNoLock(L *lua.LState, foreground bool) error {
	if foreground {
		c.foreground = true
	}
	if c.gocmd.Process != nil {
		return nil
	}

	var err error

	path := c.gocmd.Args[0]
	if c.AutoLookPath && !strings.ContainsRune(path, '/') {
		path, err = exec.LookPath(path)
		if err != nil {
			return err
		}
	}
	c.gocmd.Path = path

	err = c.setupStdio(foreground)
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
		_ = c.setupRemoteStdio(L)
		_ = c.gocmd.Wait()
		c.waitSync.Done()
	}()

	return nil
}

func (c *Cmd) ensurePrepared(L *lua.LState) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if c.gocmd.Process == nil {
		if err := c.prepareAndStartNoLock(L, true); err != nil {
			L.Pop(handleCmdExitNoLock(L, err, ExitCodeProcessCouldNotStart, c))
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
		_ = c.gocmd.Process.Kill()
	}

	L.Push(ud)
	return 1
}
