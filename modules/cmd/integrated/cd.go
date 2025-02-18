package integrated

import (
	"context"
	"os"
	"os/exec"
)

type CDCmd struct {
}

var _ Cmd = &CDCmd{}

func (c *CDCmd) RunAs(gocmd *exec.Cmd) (int, error) {
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

func (c *CDCmd) SetContext(ctx context.Context) {

}
