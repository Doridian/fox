package builtin

import (
	"context"
	"os/exec"

	"github.com/Doridian/fox/modules/loader"
)

type Cmd interface {
	RunAs(ctx context.Context, loader *loader.LuaModule, gocmd *exec.Cmd) (int, error)
}

var cmdMap = make(map[string]func() Cmd)

func Lookup(arg0 string) Cmd {
	ctor := cmdMap[arg0]
	if ctor == nil {
		return nil
	}
	return ctor()
}

func Register(arg0 string, ctor func() Cmd) {
	cmdMap[arg0] = ctor
}
