package cmd

import (
	"os"

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
	return val.PushNew(L)
}

func setStdin(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	ok, p, _ := pipe.Check(L, 2, true)
	if !ok {
		return 0
	}
	if p != nil && !p.CanRead() {
		L.ArgError(2, "pipe must be a reader")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	if c.stdin != nil && !c.stdin.IsNull() && c.stdin.CanWrite() {
		c.lock.Unlock()
		L.RaiseError("stdin piped, can't redirect")
		return 0
	}

	c.stdin = p
	c.closeStdin = doClose
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
	if c.stdin != nil && !c.stdin.IsNull() && c.stdin.CanWrite() {
		p := c.stdin
		c.lock.Unlock()
		return p.PushNew(L)
	}

	stdinPipe, err := c.gocmd.StdinPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	p := pipe.NewWritePipe(c, "stdin", stdinPipe)
	c.stdin = p
	c.closeStdin = true
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
	return val.PushNew(L)
}

func setStderr(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	ok, p, _ := pipe.Check(L, 2, true)
	if !ok {
		return 0
	}
	if p != nil && !p.CanWrite() {
		L.ArgError(2, "pipe must be a writer")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	if c.stderr != nil && !c.stderr.IsNull() && c.stderr.CanRead() {
		c.lock.Unlock()
		L.RaiseError("stderr piped, can't redirect")
		return 0
	}

	c.stderr = p
	c.closeStderr = doClose
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
	if c.stderr != nil && !c.stderr.IsNull() && c.stderr.CanRead() {
		p := c.stderr
		c.lock.Unlock()
		return p.PushNew(L)
	}

	stderrPipe, err := c.gocmd.StderrPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}

	p := pipe.NewReadPipe(c, "stderr", stderrPipe)
	c.stderr = p
	c.closeStderr = true
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
	return val.PushNew(L)
}

func setStdout(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	ok, p, _ := pipe.Check(L, 2, true)
	if !ok {
		return 0
	}
	if p != nil && !p.CanWrite() {
		L.ArgError(2, "pipe must be a writer")
		return 0
	}

	doClose := L.OptBool(3, true)

	c.lock.Lock()
	if c.stdout != nil && !c.stdout.IsNull() && c.stdout.CanRead() {
		c.lock.Unlock()
		L.RaiseError("stdout piped, can't redirect")
		return 0
	}

	c.stdout = p
	c.closeStdout = doClose
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
	if c.stdout != nil && !c.stdout.IsNull() && c.stdout.CanRead() {
		p := c.stdout
		c.lock.Unlock()
		return p.PushNew(L)
	}

	stdoutPipe, err := c.gocmd.StdoutPipe()
	if err != nil {
		c.lock.Unlock()
		L.RaiseError("%v", err)
		return 0
	}
	p := pipe.NewReadPipe(c, "stdout", stdoutPipe)
	c.stdout = p
	c.closeStdout = true
	c.lock.Unlock()
	return p.PushNew(L)
}

func (c *Cmd) setupStdio(defaultStdin bool) error {
	if c.stdout != nil {
		if c.stdout.CanWrite() {
			c.gocmd.Stdout = c.stdout.GetWriter()
		}
	} else {
		c.gocmd.Stdout = os.Stdout
	}
	if c.stderr != nil {
		if c.stderr.CanWrite() {
			c.gocmd.Stderr = c.stderr.GetWriter()
		}
	} else {
		c.gocmd.Stderr = os.Stderr
	}
	if c.stdin != nil {
		if c.stdin.CanRead() {
			c.gocmd.Stdin = c.stdin.GetReader()
		}
	} else if defaultStdin {
		c.gocmd.Stdin = os.Stdin
	} else {
		c.gocmd.Stdin = nil
	}
	return nil
}

func (c *Cmd) waitStdio() error {
	c.lock.RLock()
	stdin := c.stdin
	c.lock.RUnlock()

	if stdin == nil {
		return nil
	}
	creator := stdin.Creator()
	if creator == nil {
		return nil
	}
	cmd, ok := creator.(*Cmd)
	if !ok || cmd == nil || cmd == c {
		return nil
	}

	return cmd.ensureRan()
}

func (c *Cmd) releaseStdioNoLock() error {
	if c.stdin != nil {
		if c.closeStdin {
			c.stdin.Close()
		}
		c.stdin = nil
	}
	if c.stdout != nil {
		if c.closeStdout {
			c.stdout.Close()
		}
		c.stdout = nil
	}
	if c.stderr != nil {
		if c.closeStderr {
			c.stderr.Close()
		}
		c.stderr = nil
	}
	return nil
}
