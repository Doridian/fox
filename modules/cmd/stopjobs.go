package cmd

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/modules/cmd/integrated"
)

type StopJobsCmd struct {
}

var _ integrated.Cmd = &StopJobsCmd{}

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
	integrated.Register("stopjobs", func() integrated.Cmd { return &StopJobsCmd{} })
}
