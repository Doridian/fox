package cmd

import (
	"io"
	"os"

	luaio "github.com/Doridian/fox/modules/io"
	"github.com/Doridian/fox/modules/pipe"
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
		c.stdinCloser = nil
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	pipeIo, ud := luaio.Check(L, 2)
	if !ioCanRead(pipeIo) {
		L.ArgError(2, "pipe must be a reader")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	if c.stdin != nil && !ioIsNull(c.stdin) && ioCanWrite(c.stdin) {
		c.lock.Unlock()
		L.RaiseError("stdin piped, can't redirect")
		return 0
	}

	c.stdin = pipeIo
	if doClose {
		c.stdinCloser, _ = pipeIo.(io.Closer)
	} else {
		c.stdinCloser = nil
	}
	c.lock.Unlock()
	L.Push(ud)
	return 1
}

func acquireStdinPipe(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stdin != nil && !ioIsNull(c.stdin) && ioCanWrite(c.stdin) {
		p := c.stdin
		c.lock.Unlock()
		return luaio.PushNew(L, p)
	}

	stdinPipe, err := c.gocmd.StdinPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	p := pipe.NewPipe(c, "stdin", c.gocmd.Stdin, stdinPipe, stdinPipe)
	c.stdin = p
	c.stdinCloser, _ = c.gocmd.Stdin.(io.Closer)
	c.lock.Unlock()
	return p.PushNew(L)
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
		c.stderrCloser = nil
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	pipeIo, ud := luaio.Check(L, 2)
	if !ioCanWrite(pipeIo) {
		L.ArgError(2, "pipe must be a writer")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	if c.stderr != nil && !ioIsNull(c.stderr) && ioCanRead(c.stderr) {
		c.lock.Unlock()
		L.RaiseError("stderr piped, can't redirect")
		return 0
	}

	c.stderr = pipeIo
	if doClose {
		c.stderrCloser, _ = pipeIo.(io.Closer)
	} else {
		c.stderrCloser = nil
	}
	c.lock.Unlock()
	L.Push(ud)
	return 1
}

func acquireStderrPipe(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stderr != nil && !ioIsNull(c.stderr) && ioCanRead(c.stderr) {
		p := c.stderr
		c.lock.Unlock()
		return luaio.PushNew(L, p)
	}

	stderrPipe, err := c.gocmd.StderrPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	p := pipe.NewPipe(c, "stdout", stderrPipe, c.gocmd.Stderr, stderrPipe)
	c.stderr = p
	c.stderrCloser, _ = c.gocmd.Stderr.(io.Closer)
	c.lock.Unlock()
	return p.PushNew(L)
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
		c.stdoutCloser = nil
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	pipeIo, ud := luaio.Check(L, 2)
	if !ioCanWrite(pipeIo) {
		L.ArgError(2, "pipe must be a writer")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	if c.stdout != nil && !ioIsNull(c.stdout) && ioCanRead(c.stdout) {
		c.lock.Unlock()
		L.RaiseError("stdout piped, can't redirect")
		return 0
	}

	c.stdout = pipeIo
	if doClose {
		c.stdoutCloser, _ = pipeIo.(io.Closer)
	} else {
		c.stdoutCloser = nil
	}
	c.lock.Unlock()
	L.Push(ud)
	return 1
}

func acquireStdoutPipe(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stdout != nil && !ioIsNull(c.stdout) && ioCanRead(c.stdout) {
		p := c.stdout
		c.lock.Unlock()
		return luaio.PushNew(L, p)
	}

	stdoutPipe, err := c.gocmd.StdoutPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}
	p := pipe.NewPipe(c, "stdout", stdoutPipe, c.gocmd.Stdout, stdoutPipe)
	c.stdout = p
	c.stdoutCloser, _ = c.gocmd.Stdout.(io.Closer)
	c.lock.Unlock()
	return p.PushNew(L)
}

func (c *Cmd) setupStdio(defaultStdin bool) error {
	if c.stdout != nil {
		if !c.ioIsSelf(c.stdout) {
			c.gocmd.Stdout = c.stdout.(io.Writer)
		}
	} else {
		c.gocmd.Stdout = os.Stdout
	}
	if c.stderr != nil {
		if !c.ioIsSelf(c.stderr) {
			c.gocmd.Stderr = c.stderr.(io.Writer)
		}
	} else {
		c.gocmd.Stderr = os.Stderr
	}
	if c.stdin != nil {
		if !c.ioIsSelf(c.stdin) {
			c.gocmd.Stdin = c.stdin.(io.Reader)
		}
	} else if defaultStdin {
		c.gocmd.Stdin = os.Stdin
	} else {
		c.gocmd.Stdin = nil
	}

	return nil
}

func (c *Cmd) waitDepStdio(L *lua.LState, doAwait bool) error {
	c.lock.RLock()
	stdin := c.stdin
	c.lock.RUnlock()

	if stdin == nil {
		return nil
	}
	stdinPipe, ok := stdin.(*pipe.Pipe)
	if !ok {
		return nil
	}
	creator := stdinPipe.Creator()
	if creator == nil {
		return nil
	}

	cmd, ok := creator.(*Cmd)
	if !ok || cmd == nil || cmd == c {
		return nil
	}

	return cmd.ensureRan(L, doAwait)
}

func (c *Cmd) releaseStdioNoLock() error {
	if c.stdinCloser != nil {
		_ = ioClose(c.stdinCloser)
		c.stdinCloser = nil
	}
	if c.stdoutCloser != nil {
		_ = ioClose(c.stdoutCloser)
		c.stdoutCloser = nil
	}
	if c.stderrCloser != nil {
		_ = ioClose(c.stderrCloser)
		c.stderrCloser = nil
	}

	return nil
}
