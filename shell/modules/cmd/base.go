package cmd

import (
	"os/exec"
	"sync"

	"github.com/Doridian/fox/shell/modules/pipe"
	lua "github.com/yuin/gopher-lua"
)

type Cmd struct {
	stdout      *pipe.Pipe
	closeStdout bool
	stderr      *pipe.Pipe
	closeStderr bool
	stdin       *pipe.Pipe
	closeStdin  bool

	lock             sync.RWMutex
	gocmd            *exec.Cmd
	ErrorPropagation bool
}

func newCmd(L *lua.LState) int {
	return pushCmd(L, &Cmd{
		gocmd:            &exec.Cmd{},
		ErrorPropagation: false,
	})
}
