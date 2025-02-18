package integrated

import (
	"context"
	"os"
	"os/exec"
)

type PwdCmd struct {
}

var _ Cmd = &PwdCmd{}

func (c *PwdCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 1, err
	}
	gocmd.Stdout.Write([]byte(wd))
	gocmd.Stdout.Write([]byte("\n"))
	return 0, nil
}

func (c *PwdCmd) SetContext(ctx context.Context) {

}
