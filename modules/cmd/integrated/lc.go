package integrated

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
		_, _ = gocmd.Stderr.Write([]byte("lc: missing command name\n"))
		return 1, nil
	}

	subShell := shell.New(c.ctx)
	subShell.ShowErrors = false
	defer subShell.Close()

	err := subShell.Init(gocmd.Args[1:], false)
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
