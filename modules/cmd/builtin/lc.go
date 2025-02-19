package builtin

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/shell"
)

type LCCmd struct {
	ctx context.Context
}

var _ Cmd = &LCCmd{}

func (c *LCCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		_, _ = gocmd.Stderr.Write([]byte("missing command name\n"))
		return 1, nil
	}

	subArgs := gocmd.Args[1:]

	if gocmd.Args[1] == "-p" {
		if len(gocmd.Args) < 3 {
			_, _ = gocmd.Stderr.Write([]byte("missing command name\n"))
			return 1, nil
		}

		var currentShell *shell.Shell // TODO: How to get the current shell?
		subArgs = currentShell.GetArgs()
	}

	subShell := shell.New(c.ctx)
	subShell.ShowErrors = false
	defer subShell.Close()

	err := subShell.Init(subArgs, false)
	if err != nil {
		return 1, err
	}
	subShell.SetStdio(gocmd.Stdin, gocmd.Stdout, gocmd.Stderr)
	err = subShell.RunCommand(gocmd.Args[1])
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func (c *LCCmd) SetContext(ctx context.Context) {
	c.ctx = ctx
}

func init() {
	Register("lc", func() Cmd { return &LCCmd{} })
}
