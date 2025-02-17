package integrated

import (
	"context"
	"os/exec"
)

type EchoCmd struct {
}

func (c *EchoCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	for i, arg := range gocmd.Args[1:] {
		if i > 0 {
			gocmd.Stdout.Write([]byte(" "))
		}
		gocmd.Stdout.Write([]byte(arg))
	}
	gocmd.Stdout.Write([]byte("\n"))
	return 0, nil
}

func (c *EchoCmd) SetContext(ctx context.Context) {

}

var _ Cmd = &EchoCmd{}
