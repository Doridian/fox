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
		c.stdinCloser = nil
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
		c.stdinCloser, _ = ioL.(io.Closer)
	} else {
		c.stdinCloser = nil
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

	c.stdin = nil
	c.stdinPipe = stdinPipe
	c.stdinCloser, _ = c.gocmd.Stdin.(io.Closer)
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
		c.stderrCloser = nil
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
		c.stderrCloser, _ = ioL.(io.Closer)
	} else {
		c.stderrCloser = nil
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

	c.stderr = nil
	c.stderrPipe = stderrPipe
	c.stderrCloser, _ = c.gocmd.Stderr.(io.Closer)
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
		c.stdoutCloser = nil
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
		c.stdoutCloser, _ = ioL.(io.Closer)
	} else {
		c.stdoutCloser = nil
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

	c.stdout = nil
	c.stdoutPipe = stdoutPipe
	c.stdoutCloser, _ = c.gocmd.Stdout.(io.Closer)
	c.lock.Unlock()
	return luaio.PushNew(L, stdoutPipe)
}

func (c *Cmd) setupStdio(defaultStdin bool) error {
	if c.stdout != nil {
		c.gocmd.Stdout = c.stdout
	} else if c.stdoutPipe == nil {
		c.gocmd.Stdout = shell.StdoutFor(c.mod.loader)
	}
	if c.stderr != nil {
		c.gocmd.Stderr = c.stderr
	} else if c.stderrPipe == nil {
		c.gocmd.Stderr = shell.StderrFor(c.mod.loader)
	}
	if c.stdin != nil {
		c.gocmd.Stdin = c.stdin
	} else if c.stdinPipe == nil {
		if defaultStdin {
			c.gocmd.Stdin = shell.StdinFor(c.mod.loader)
		} else {
			c.gocmd.Stdin = nil
		}
	}

	return nil
}

func (c *Cmd) releaseStdioNoLock() error {
	if c.stdinCloser != nil {
		_ = c.stderrCloser.Close()
		c.stdinCloser = nil
	}
	if c.stdoutCloser != nil {
		_ = c.stdoutCloser.Close()
		c.stdoutCloser = nil
	}
	if c.stderrCloser != nil {
		_ = c.stderrCloser.Close()
		c.stderrCloser = nil
	}

	return nil
}
