package cmd

import (
	"os"

	"github.com/Doridian/fox/shell/modules/pipe"
	lua "github.com/yuin/gopher-lua"
)

func getSetStdin(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		ok, p, _ := pipe.CheckPipe(L, 2, true)
		if !ok {
			return 0
		}
		if p != nil && !p.CanRead() {
			L.ArgError(2, "pipe must be a reader")
			return 0
		}

		c.lock.Lock()
		if c.stdin != nil && !c.stdin.IsNull() && c.stdin.CanWrite() {
			c.lock.Unlock()
			L.Error(lua.LString("stdin piped, can't redirect"), 0)
			return 0
		}

		c.stdin = p
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	val := c.stdin
	c.lock.RUnlock()
	return val.Push(L)
}

func getStdinPipe(L *lua.LState) int {
	c, _ := checkCmd(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stdin != nil && !c.stdin.IsNull() && c.stdin.CanWrite() {
		p := c.stdin
		c.lock.Unlock()
		return p.Push(L)
	}

	stdinPipe, err := c.gocmd.StdinPipe()
	if err != nil {
		c.lock.Unlock()
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	p := pipe.NewWritePipe(c, stdinPipe)
	c.stdin = p
	c.lock.Unlock()
	return p.Push(L)
}

func getSetStderr(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		ok, p, _ := pipe.CheckPipe(L, 2, true)
		if !ok {
			return 0
		}
		if p != nil && !p.CanWrite() {
			L.ArgError(2, "pipe must be a writer")
			return 0
		}

		c.lock.Lock()
		if c.stderr != nil && !c.stderr.IsNull() && c.stderr.CanRead() {
			c.lock.Unlock()
			L.Error(lua.LString("stderr piped, can't redirect"), 0)
			return 0
		}

		c.stderr = p
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	val := c.stderr
	c.lock.RUnlock()
	return val.Push(L)
}

func getStderrPipe(L *lua.LState) int {
	c, _ := checkCmd(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stderr != nil && !c.stderr.IsNull() && c.stderr.CanRead() {
		p := c.stderr
		c.lock.Unlock()
		return p.Push(L)
	}

	stderrPipe, err := c.gocmd.StderrPipe()
	if err != nil {
		c.lock.Unlock()
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	p := pipe.NewReadPipe(c, stderrPipe)
	c.stderr = p
	c.lock.Unlock()
	return p.Push(L)
}

func getSetStdout(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		ok, p, _ := pipe.CheckPipe(L, 2, true)
		if !ok {
			return 0
		}
		if p != nil && !p.CanWrite() {
			L.ArgError(2, "pipe must be a writer")
			return 0
		}

		c.lock.Lock()
		if c.stdout != nil && !c.stdout.IsNull() && c.stdout.CanRead() {
			c.lock.Unlock()
			L.Error(lua.LString("stdout piped, can't redirect"), 0)
			return 0
		}

		c.stdout = p
		c.lock.Unlock()
		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	val := c.stdout
	c.lock.RUnlock()
	return val.Push(L)
}

func getStdoutPipe(L *lua.LState) int {
	c, _ := checkCmd(L, 1)
	if c == nil {
		return 0
	}

	c.lock.Lock()
	if c.stdout != nil && !c.stdout.IsNull() && c.stdout.CanRead() {
		p := c.stdout
		c.lock.Unlock()
		return p.Push(L)
	}

	stdoutPipe, err := c.gocmd.StdoutPipe()
	if err != nil {
		c.lock.Unlock()
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	p := pipe.NewReadPipe(c, stdoutPipe)
	c.stdout = p
	c.lock.Unlock()
	return p.Push(L)
}

func (c *Cmd) setupStdio() error {
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
	} else {
		c.gocmd.Stdin = os.Stdin
	}
	return nil
}

func (c *Cmd) waitStdio() error {
	if c.stdin != nil {
		creator := c.stdin.Creator()
		if creator != nil {
			cmd, ok := creator.(*Cmd)
			if ok && cmd != nil {
				return cmd.prepareAndRun()
			}
		}
	}
	return nil
}

func (c *Cmd) releaseStdio() error {
	if c.stdin != nil {
		c.stdin.Close()
		c.stdin = nil
	}
	if c.stdout != nil {
		c.stdout.Close()
		c.stdout = nil
	}
	if c.stderr != nil {
		c.stderr.Close()
		c.stderr = nil
	}
	return nil
}
