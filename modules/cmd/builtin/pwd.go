package builtin

import (
	"context"
	"os"
	"os/exec"

	"github.com/Doridian/fox/modules/loader"
)

type PwdCmd struct {
}

var _ Cmd = &PwdCmd{}

func (c *PwdCmd) RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error) {
	wd, err := os.Getwd()
	if err != nil {
		return 1, err
	}
	_, _ = gocmd.Stdout.Write([]byte(wd))
	_, _ = gocmd.Stdout.Write([]byte("\n"))
	return 0, nil
}

func (c *PwdCmd) SetContext(ctx context.Context) {

}

func init() {
	Register("pwd", func() Cmd { return &PwdCmd{} })
}
