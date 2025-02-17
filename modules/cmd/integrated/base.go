package integrated

import (
	"context"
	"os/exec"
)

type Cmd interface {
	SetContext(ctx context.Context)
	RunAs(gocmd *exec.Cmd) (int, error)
}

func LookupCmd(arg0 string) Cmd {
	switch arg0 {
	case "echo":
		return &EchoCmd{}
	}

	return nil
}
