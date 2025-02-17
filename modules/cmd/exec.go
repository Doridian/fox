package cmd

import (
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func handleCmdExitNoLock(L *lua.LState, nonExitError error, exitCode int, c *Cmd) int {
	_ = c.releaseStdioNoLock()
	if L == nil {
		return 0
	}

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

func (c *Cmd) doWaitCmdNoLock(L *lua.LState) {
	c.awaited = true
	_ = c.waitDepStdio(L, true)
	c.waitSync.Wait()
}

func doWaitCmdNoLock(L *lua.LState, c *Cmd) int {
	c.doWaitCmdNoLock(L)
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

	err := c.prepareAndStartNoLock(true)
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

	err := c.prepareAndStartNoLock(false)
	if err != nil {
		return handleCmdExitNoLock(L, err, ExitCodeProcessCouldNotStart, c)
	}
	return 0
}

func (c *Cmd) prepareAndStartNoLock(foreground bool) error {
	c.startLock.Lock()
	defer c.startLock.Unlock()

	if foreground {
		c.foreground = true
	}
	if c.gocmd.Process != nil || c.gocmd.ProcessState != nil {
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
		_ = c.waitDepStdio(nil, false)
		_ = c.gocmd.Wait()
		c.waitSync.Done()
	}()

	return nil
}

func (c *Cmd) ensureRan(L *lua.LState, doAwait bool) error {
	c.lock.RLock()
	defer c.lock.RUnlock()

	if err := c.prepareAndStartNoLock(true); err != nil {
		retC := handleCmdExitNoLock(L, err, ExitCodeProcessCouldNotStart, c)
		if retC > 0 {
			L.Pop(retC)
		}
		return err
	}

	if doAwait {
		c.doWaitCmdNoLock(L)
	}
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
