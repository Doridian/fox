package cmd

import (
	"io"

	luaio "github.com/Doridian/fox/modules/io"
	"github.com/Doridian/fox/shell"
	lua "github.com/yuin/gopher-lua"
)

func getStdin(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := c.stdin
	c.lock.RUnlock()
	return luaio.PushNew(L, val)
}

func setStdin(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	if L.Get(2) == lua.LNil {
		c.lock.Lock()
		c.stdin = nil
		c.stdinPipe = nil
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	ioL, ud := luaio.Check(L, 2)
	r, ok := ioL.(io.Reader)
	if !ok {
		L.ArgError(2, "pipe must be a reader")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.stdin = r
	c.stdinPipe = nil
	if doClose {
		c.addCloserNoLock(ioL)
	}
	L.Push(ud)
	return 1
}

func acquireStdinPipe(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stdinPipe != nil {
		p := c.stdinPipe
		c.lock.Unlock()
		return luaio.PushNew(L, p)
	}

	stdinPipe, err := c.gocmd.StdinPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	c.stdin = c.gocmd.Stdin
	c.stdinPipe = stdinPipe
	c.addCloserNoLock(c.gocmd.Stdin)
	c.lock.Unlock()
	return luaio.PushNew(L, stdinPipe)
}

func getStderr(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := c.stderr
	c.lock.RUnlock()
	return luaio.PushNew(L, val)
}

func setStderr(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	if L.Get(2) == lua.LNil {
		c.lock.Lock()
		c.stderr = nil
		c.stderrPipe = nil
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	ioL, ud := luaio.Check(L, 2)
	w, ok := ioL.(io.Writer)
	if !ok {
		L.ArgError(2, "pipe must be a writer")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.stderr = w
	c.stderrPipe = nil
	if doClose {
		c.addCloserNoLock(ioL)
	}
	L.Push(ud)
	return 1
}

func acquireStderrPipe(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stderrPipe != nil {
		p := c.stderrPipe
		c.lock.Unlock()
		return luaio.PushNew(L, p)
	}

	stderrPipe, err := c.gocmd.StderrPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	c.stderr = c.gocmd.Stderr
	c.stderrPipe = stderrPipe
	c.addCloserNoLock(c.gocmd.Stderr)
	c.lock.Unlock()
	return luaio.PushNew(L, stderrPipe)
}

func getStdout(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := c.stdout
	c.lock.RUnlock()
	return luaio.PushNew(L, val)
}

func setStdout(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	if L.Get(2) == lua.LNil {
		c.lock.Lock()
		c.stdout = nil
		c.stdoutPipe = nil
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	ioL, ud := luaio.Check(L, 2)
	w, ok := ioL.(io.Writer)
	if !ok {
		L.ArgError(2, "pipe must be a writer")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	defer c.lock.Unlock()

	c.stdout = w
	c.stdoutPipe = nil
	if doClose {
		c.addCloserNoLock(ioL)
	}
	L.Push(ud)
	return 1
}

func acquireStdoutPipe(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stdoutPipe != nil {
		p := c.stdoutPipe
		c.lock.Unlock()
		return luaio.PushNew(L, p)
	}

	stdoutPipe, err := c.gocmd.StdoutPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	c.stdout = c.gocmd.Stdout
	c.stdoutPipe = stdoutPipe
	c.addCloserNoLock(c.gocmd.Stdout)
	c.lock.Unlock()
	return luaio.PushNew(L, stdoutPipe)
}

func (c *Cmd) setupStdio(defaultStdin bool) error {
	if c.stdoutPipe == nil {
		if c.stdout != nil {
			c.gocmd.Stdout = c.stdout
		} else {
			c.gocmd.Stdout = shell.StdoutFor(c.mod.loader)
		}
	}
	if c.stderrPipe == nil {
		if c.stderr != nil {
			c.gocmd.Stderr = c.stderr
		} else {
			c.gocmd.Stderr = shell.StderrFor(c.mod.loader)
		}
	}
	if c.stdinPipe == nil {
		if c.stdin != nil {
			c.gocmd.Stdin = c.stdin
		} else if defaultStdin {
			c.gocmd.Stdin = shell.StdinFor(c.mod.loader)
		} else {
			c.gocmd.Stdin = nil
		}
	}

	return nil
}

func (c *Cmd) addCloserNoLock(closerRaw interface{}) {
	closer, ok := closerRaw.(io.Closer)
	if !ok {
		return
	}
	c.closeQueue = append(c.closeQueue, closer)
}

func (c *Cmd) releaseStdioNoLock() error {
	for _, closer := range c.closeQueue {
		closer.Close()
	}
	c.closeQueue = make([]io.Closer, 0)

	return nil
}
