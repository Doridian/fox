package integrated

import (
	"context"
	"os/exec"
)

type EchoCmd struct {
}

func (e *EchoCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	for i, arg := range gocmd.Args[1:] {
		if i > 0 {
			gocmd.Stdout.Write([]byte(" "))
		}
		gocmd.Stdout.Write([]byte(arg))
	}
	return 0, nil
}

func (e *EchoCmd) SetContext(ctx context.Context) {

}

var _ Cmd = &EchoCmd{}
