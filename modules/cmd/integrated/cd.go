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
		gocmd.Stderr.Write([]byte("cd: missing dir\n"))
		return 1, nil
	}

	os.Chdir(gocmd.Args[1])
	return 0, nil
}

func (c *CDCmd) SetContext(ctx context.Context) {

}
