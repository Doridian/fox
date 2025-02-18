package integrated

import (
	"context"
	"os"
	"os/exec"
	"strconv"
)

type ExitCmd struct {
}

var _ Cmd = &ExitCmd{}

func (c *ExitCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	if len(gocmd.Args) < 2 {
		os.Exit(0)
		return 0, nil
	}

	code, err := strconv.ParseInt(gocmd.Args[1], 10, 32)
	if err != nil {
		gocmd.Stderr.Write([]byte(err.Error()))
		gocmd.Stderr.Write([]byte("\n"))
		code = 1
	}
	os.Exit(int(code))
	return int(code), nil
}

func (c *ExitCmd) SetContext(ctx context.Context) {

}
