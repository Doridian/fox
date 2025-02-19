package cmd

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/modules/cmd/builtin"
	"github.com/Doridian/fox/modules/loader"
)

type StopAllCmd struct {
}

var _ builtin.Cmd = &StopAllCmd{}

func (c *StopAllCmd) RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error) {
	cmdRegLock.Lock()
	for cmd := range allCmds {
		cmd.Stop()
	}
	cmdRegLock.Unlock()
	return 0, nil
}

func init() {
	builtin.Register("stopall", func() builtin.Cmd { return &StopAllCmd{} })
}
