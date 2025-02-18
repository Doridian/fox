package integrated

import (
	"context"
	"io"
	"os/exec"
)

type PSICmd struct {
}

func (c *PSICmd) RunAs(gocmd *exec.Cmd) (int, error) {
	varB, err := io.ReadAll(gocmd.Stdin)
	if err != nil {
		return 1, err
	}
	_, _ = gocmd.Stdout.Write(varB)

	return 0, nil
}

func (c *PSICmd) SetContext(ctx context.Context) {

}

var _ Cmd = &PSICmd{}
