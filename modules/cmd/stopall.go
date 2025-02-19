package cmd

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/modules/cmd/builtin"
)

type StopAllCmd struct {
}

var _ builtin.Cmd = &StopAllCmd{}

func (c *StopAllCmd) RunAs(gocmd *exec.Cmd) (int, error) {
	cmdRegLock.Lock()
	for cmd := range allCmds {
		cmd.Stop()
	}
	cmdRegLock.Unlock()
	return 0, nil
}

func (c *StopAllCmd) SetContext(ctx context.Context) {

}

func init() {
	builtin.Register("stopall", func() builtin.Cmd { return &StopAllCmd{} })
}
