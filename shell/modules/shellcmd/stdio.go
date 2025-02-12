package shellcmd

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

func getSetStdin(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		ok, pipe, _ := checkPipe(L, 2, true)
		if !ok {
			return 0
		}
		if pipe != nil && pipe.rc == nil {
			L.ArgError(2, "stdin must be a reader")
			return 0
		}

		if c.stdin != nil && c.stdin.wc != nil {
			L.Error(lua.LString("stdin piped, can't redirect"), 0)
			return 0
		}

		c.stdin = pipe
		L.Push(ud)
		return 1
	}

	return pushShellPipe(L, c.stdin)
}

func getStdinPipe(L *lua.LState) int {
	c, _ := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if c.stdin != nil && c.stdin.wc != nil {
		return pushShellPipe(L, c.stdin)
	}

	pipe, err := c.gocmd.StdinPipe()
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	c.stdin = &Pipe{cmd: c, wc: pipe, forwardClose: true}
	return pushShellPipe(L, c.stdin)
}

func getSetStderr(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		ok, pipe, _ := checkPipe(L, 2, true)
		if !ok {
			return 0
		}
		if pipe != nil && pipe.wc == nil {
			L.ArgError(2, "stderr must be a writer")
			return 0
		}

		if c.stderr != nil && c.stderr.rc != nil {
			L.Error(lua.LString("stderr piped, can't redirect"), 0)
			return 0
		}

		c.stderr = pipe
		L.Push(ud)
		return 1
	}

	return pushShellPipe(L, c.stderr)
}

func getStderrPipe(L *lua.LState) int {
	c, _ := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if c.stderr != nil && c.stderr.rc != nil {
		return pushShellPipe(L, c.stderr)
	}

	pipe, err := c.gocmd.StderrPipe()
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	c.stderr = &Pipe{cmd: c, rc: pipe, forwardClose: true}
	return pushShellPipe(L, c.stderr)
}

func getSetStdout(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		ok, pipe, _ := checkPipe(L, 2, true)
		if !ok {
			return 0
		}
		if pipe != nil && pipe.wc == nil {
			L.ArgError(2, "stdout must be a writer")
			return 0
		}

		if c.stdout != nil && c.stdout.rc != nil {
			L.Error(lua.LString("stdout piped, can't redirect"), 0)
			return 0
		}

		c.stdout = pipe
		L.Push(ud)
		return 1
	}

	return pushShellPipe(L, c.stdout)
}

func getStdoutPipe(L *lua.LState) int {
	c, _ := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if c.stdout != nil && c.stdout.rc != nil {
		return pushShellPipe(L, c.stdout)
	}

	pipe, err := c.gocmd.StdoutPipe()
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	c.stdout = &Pipe{cmd: c, rc: pipe, forwardClose: true}
	return pushShellPipe(L, c.stdout)
}

func (c *ShellCmd) setupStdio() error {
	if c.stdout != nil {
		if c.stdout.wc != nil || c.stdout.isNull {
			c.gocmd.Stdout = c.stdout.wc
		}
	} else {
		c.gocmd.Stdout = os.Stdout
	}
	if c.stderr != nil {
		if c.stderr.wc != nil || c.stderr.isNull {
			c.gocmd.Stderr = c.stderr.wc
		}
	} else {
		c.gocmd.Stderr = os.Stderr
	}
	if c.stdin != nil {
		if c.stdin.rc != nil || c.stdin.isNull {
			c.gocmd.Stdin = c.stdin.rc
		}
	} else {
		c.gocmd.Stdin = os.Stdin
	}
	return nil
}

func (c *ShellCmd) waitStdio() error {
	if c.stdin != nil {
		defer c.stdin.Close()
		if c.stdin.cmd != nil {
			return c.stdin.cmd.prepareAndRun()
		}
	}
	return nil
}

func (c *ShellCmd) releaseStdio() error {
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
