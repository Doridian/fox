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
		if c.stdinPipe != nil {
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
	if c.stdinPipe != nil {
		p := c.stdinPipe
		c.lock.Unlock()
		return p.Push(L)
	}

	cmdPipe, err := c.gocmd.StdinPipe()
	if err != nil {
		c.lock.Unlock()
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	p := pipe.NewWritePipe(c, cmdPipe)
	c.stdinPipe = p
	c.stdin = nil
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
		if c.stderrPipe != nil {
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
	if c.stderrPipe != nil {
		p := c.stderrPipe
		c.lock.Unlock()
		return p.Push(L)
	}

	cmdPipe, err := c.gocmd.StderrPipe()
	if err != nil {
		c.lock.Unlock()
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	p := pipe.NewReadPipe(c, cmdPipe)
	c.stderrPipe = p
	c.stderr = nil
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
		if c.stdoutPipe != nil {
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
	if c.stdoutPipe != nil {
		p := c.stdoutPipe
		c.lock.Unlock()
		return p.Push(L)
	}

	cmdPipe, err := c.gocmd.StdoutPipe()
	if err != nil {
		c.lock.Unlock()
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	p := pipe.NewReadPipe(c, cmdPipe)
	c.stdoutPipe = p
	c.stdout = nil
	c.lock.Unlock()
	return p.Push(L)
}

func (c *Cmd) setupStdio() error {
	if c.stdout != nil {
		c.gocmd.Stdout = c.stdout.GetWriter()
	} else if c.stdoutPipe == nil {
		c.gocmd.Stdout = os.Stdout
	}
	if c.stderr != nil {
		c.gocmd.Stderr = c.stderr.GetWriter()
	} else if c.stderrPipe == nil {
		c.gocmd.Stderr = os.Stderr
	}
	if c.stdin != nil {
		c.gocmd.Stdin = c.stdin.GetReader()
	} else if c.stdinPipe == nil {
		c.gocmd.Stdin = os.Stdin
	}
	return nil
}

func (c *Cmd) waitStdio() error {
	if c.stdin != nil {
		creator := c.stdin.Creator()
		if creator != nil {
			cmd, ok := creator.(*Cmd)
			if ok && cmd != nil && cmd != c {
				return cmd.ensureRan()
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
