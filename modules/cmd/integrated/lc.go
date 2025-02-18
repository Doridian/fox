package integrated

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/shell"
)

type LCCmd struct {
	ctx context.Context
}

func (c *LCCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	subShell := shell.New(c.ctx)
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

var _ Cmd = &LCCmd{}
