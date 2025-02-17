package cmd

import (
	"context"
	"os/exec"
	"strings"

	"github.com/Doridian/fox/modules/cmd/integrated"
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
	exitCode := c.iExit
	if c.gocmd.ProcessState != nil {
		exitCode = c.gocmd.ProcessState.ExitCode()
	}
	return handleCmdExitNoLock(L, nil, exitCode, c)
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
	if c.gocmd.Process != nil || c.gocmd.ProcessState != nil || c.iCtx != nil {
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

	c.iCmd = integrated.LookupCmd(c.gocmd.Args[0])
	if c.iCmd != nil {
		c.iCtx, c.iCancel = context.WithCancel(context.Background())
		c.iCmd.SetContext(c.iCtx)
	} else {
		err = c.gocmd.Start()
		if err != nil {
			return err
		}
	}

	c.waitSync.Add(1)
	c.mod.addCmd(c)
	go func() {
		_ = c.waitDepStdio(nil, false)
		if c.iCmd != nil {
			code, err := c.iCmd.RunAs(c.gocmd)
			c.iExit = code
			handleCmdExitNoLock(nil, err, c.iExit, c)
		} else {
			_ = c.gocmd.Wait()
		}
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

	c.startLock.Lock()
	defer c.startLock.Unlock()

	if c.gocmd.Process != nil {
		_ = c.gocmd.Process.Kill()
	}
	if c.iCancel != nil {
		c.iCancel()
	}

	L.Push(ud)
	return 1
}
