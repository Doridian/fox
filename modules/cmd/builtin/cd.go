package builtin

import (
	"context"
	"os"
	"os/exec"

	"github.com/Doridian/fox/modules/loader"
)

type CDCmd struct {
}

var _ Cmd = &CDCmd{}

func (c *CDCmd) RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		_, _ = gocmd.Stderr.Write([]byte("missing dir\n"))
		return 1, nil
	}

	err := os.Chdir(gocmd.Args[1])
	if err != nil {
		_, _ = gocmd.Stderr.Write([]byte(err.Error()))
		_, _ = gocmd.Stderr.Write([]byte("\n"))
		return 1, nil
	}
	return 0, nil
}

func init() {
	Register("cd", func() Cmd { return &CDCmd{} })
}
