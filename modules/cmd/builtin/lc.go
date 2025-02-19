package builtin

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
)

type LCCmd struct {
}

var _ Cmd = &LCCmd{}

func (c *LCCmd) RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		_, _ = gocmd.Stderr.Write([]byte("missing command name\n"))
		return 1, nil
	}

	subArgs := gocmd.Args[1:]

	if subArgs[0] == "-p" {
		if len(subArgs) < 2 {
			_, _ = gocmd.Stderr.Write([]byte("missing command name\n"))
			return 1, nil
		}

		subArgs = subArgs[1:]
		currentShell := loader.GetModule(shell.LuaName).(*shell.Shell)
		subArgs = append(subArgs, currentShell.GetArgs()...)
	}

	subShell := shell.New(ctx)
	subShell.ShowErrors = false
	defer subShell.Close()

	err := subShell.Init(subArgs, false)
	if err != nil {
		return 1, err
	}
	subShell.SetStdio(gocmd.Stdin, gocmd.Stdout, gocmd.Stderr)
	err = subShell.RunCommand(subArgs[0])
	if err != nil {
		return 1, err
	}
	return 0, nil
}

func init() {
	Register("lc", func() Cmd { return &LCCmd{} })
}
