package shellcmd

import (
	"io"
	"os"

	lua "github.com/yuin/gopher-lua"
)

func getSetStdout(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		c2, _ := checkShellCmd(L, 2)
		if c2 == nil {
			return 0
		}
		c.Stdout = c2
		L.Push(ud)
		return 1
	}

	return pushShellCmd(L, c.Stdout)
}

func getSetStderr(L *lua.LState) int {
	c, ud := checkShellCmd(L, 1)
	if c == nil {
		return 0
	}

	if L.GetTop() >= 2 {
		c2, _ := checkShellCmd(L, 2)
		if c2 == nil {
			return 0
		}
		c.Stderr = c2
		L.Push(ud)
		return 1
	}

	return pushShellCmd(L, c.Stderr)
}

func (c *ShellCmd) setupPipes() error {
	c.Gocmd.Stderr = os.Stderr
	c.Gocmd.Stdout = os.Stdout
	c.Gocmd.Stdin = os.Stdin

	var stdoutStdinPipe io.WriteCloser
	var err error
	if c.Stdout != nil {
		stdoutStdinPipe, err = c.Stdout.Gocmd.StdinPipe()
		if err != nil {
			return err
		}
		c.Gocmd.Stdout = stdoutStdinPipe
	}

	if c.Stderr != nil {
		if c.Stdout == c.Stderr {
			c.Gocmd.Stderr = stdoutStdinPipe
		} else {
			stderrStdinPipe, err := c.Stderr.Gocmd.StdinPipe()
			if err != nil {
				return err
			}
			c.Gocmd.Stderr = stderrStdinPipe
		}
	}

	return nil
}

func (c *ShellCmd) waitPipes() error {
	var errStdout, errStderr error
	if c.Stdout != nil {
		errStdout = c.Stdout.prepareAndWait()
	}

	if c.Stderr != nil && c.Stdout != c.Stderr {
		errStderr = c.Stderr.prepareAndWait()
	}

	if errStdout != nil {
		return errStdout
	}
	return errStderr
}
