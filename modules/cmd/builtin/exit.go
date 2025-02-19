package builtin

import (
	"context"
	"os"
	"os/exec"
	"strconv"

	"github.com/Doridian/fox/modules/loader"
)

type ExitCmd struct {
}

var _ Cmd = &ExitCmd{}

func (c *ExitCmd) RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		os.Exit(0)
		return 0, nil
	}

	code, err := strconv.ParseInt(gocmd.Args[1], 10, 32)
	if err != nil {
		_, _ = gocmd.Stderr.Write([]byte(err.Error()))
		_, _ = gocmd.Stderr.Write([]byte("\n"))
		code = 1
	}
	os.Exit(int(code))
	return int(code), nil
}

func init() {
	Register("exit", func() Cmd { return &ExitCmd{} })
}
