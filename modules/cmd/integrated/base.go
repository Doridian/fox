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
	case "export":
		return &ExportCmd{}
	case "set":
		return &SetCmd{}
	case "lc":
		return &LCCmd{}
	case "cd":
		return &CDCmd{}
	case "pwd":
		return &PwdCmd{}
	}

	return nil
}
