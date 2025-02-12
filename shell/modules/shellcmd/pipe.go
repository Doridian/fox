package shellcmd

import (
	"io"
	"os"

	lua "github.com/yuin/gopher-lua"
)

const luaShellPipeType = "shell/modules/shellcmd/ShellPipe"

type ShellPipe struct {
	cmd *ShellCmd
	rc  io.ReadCloser
	wc  io.WriteCloser
}

func (p *ShellPipe) Close() {
	if p.rc != nil {
		p.rc.Close()
	}
	if p.wc != nil {
		p.wc.Close()
	}
}

func newNullShellPipe(L *lua.LState) int {
	return pushShellPipe(L, &ShellPipe{})
}

func getSetStdin(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		pipe, _ := checkShellPipe(L, 2)
		if pipe.rc == nil {
			L.Error(lua.LString("stdin pipe must be a reader"), 0)
			return 0
		}

		if c.stdin != nil && c.stdin.wc != nil {
			c.stdin.wc.Close()
		}
		c.stdin = pipe

		if pipe == nil {
			c.gocmd.Stdin = os.Stdin
		} else {
			c.gocmd.Stdin = pipe.rc
		}
		L.Push(ud)
		return 1
	}

	if c.stdin == nil {
		c.gocmd.Stdin = nil
		pipe, err := c.gocmd.StdinPipe()
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		c.stdin = &ShellPipe{cmd: c, wc: pipe}
	}

	return pushShellPipe(L, c.stdin)
}

func getSetStderr(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		pipe, _ := checkShellPipe(L, 2)
		if pipe.wc == nil {
			L.Error(lua.LString("stderr pipe must be a writer"), 0)
			return 0
		}

		if c.stderr != nil && c.stderr.rc != nil {
			c.stderr.rc.Close()
		}
		c.stderr = pipe

		if pipe == nil {
			c.gocmd.Stderr = os.Stderr
		} else {
			c.gocmd.Stderr = pipe.wc
		}
		L.Push(ud)
		return 1
	}

	if c.stderr == nil {
		c.gocmd.Stderr = nil
		pipe, err := c.gocmd.StderrPipe()
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		c.stderr = &ShellPipe{cmd: c, rc: pipe}
	}

	return pushShellPipe(L, c.stderr)
}

func getSetStdout(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		pipe, _ := checkShellPipe(L, 2)
		if pipe.wc == nil {
			L.Error(lua.LString("stdout pipe must be a writer"), 0)
			return 0
		}

		if c.stdout != nil && c.stdout.rc != nil {
			c.stdout.rc.Close()
		}
		c.stdout = pipe

		if pipe == nil {
			c.gocmd.Stdout = os.Stdout
		} else {
			c.gocmd.Stdout = pipe.wc
		}
		L.Push(ud)
		return 1
	}

	if c.stdout == nil {
		c.gocmd.Stdout = nil
		pipe, err := c.gocmd.StdoutPipe()
		if err != nil {
			L.Error(lua.LString(err.Error()), 0)
			return 0
		}
		c.stdout = &ShellPipe{cmd: c, rc: pipe}
	}

	return pushShellPipe(L, c.stdout)
}

func (c *ShellCmd) setupPipes() error {
	return nil
}

func (c *ShellCmd) waitPipes() error {
	if c.stdin != nil {
		defer c.stdin.Close()
		return c.stdin.cmd.prepareAndRun()
	}
	return nil
}

func (c *ShellCmd) releasePipes() error {
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
