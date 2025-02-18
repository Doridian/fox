package integrated

import (
	"context"
	"os/exec"
)

type EchoCmd struct {
}

var _ Cmd = &EchoCmd{}

func (c *EchoCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		return 0, nil
	}

	for i, arg := range gocmd.Args[1:] {
		if i > 0 {
			_, _ = gocmd.Stdout.Write([]byte(" "))
		}
		_, _ = gocmd.Stdout.Write([]byte(arg))
	}
	_, _ = gocmd.Stdout.Write([]byte("\n"))
	return 0, nil
}

func (c *EchoCmd) SetContext(ctx context.Context) {

}

func init() {
	Register("echo", func() Cmd { return &EchoCmd{} })
}
