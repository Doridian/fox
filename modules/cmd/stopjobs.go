package cmd

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/modules/cmd/builtin"
)

type StopJobsCmd struct {
}

var _ builtin.Cmd = &StopJobsCmd{}

func (c *StopJobsCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	cmdRegLock.Lock()
	for cmd := range allCmds {
		cmd.Stop()
	}
	cmdRegLock.Unlock()
	return 0, nil
}

func (c *StopJobsCmd) SetContext(ctx context.Context) {

}

func init() {
	builtin.Register("stopjobs", func() builtin.Cmd { return &StopJobsCmd{} })
}
