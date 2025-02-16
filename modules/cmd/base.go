package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/Doridian/fox/modules/pipe"
)

const (
	ExitCodeLuaError             = -10001
	ExitCodeProcessCouldNotStart = -10002
)

type Cmd struct {
	stdout      *pipe.Pipe
	closeStdout bool
	stderr      *pipe.Pipe
	closeStderr bool
	stdin       *pipe.Pipe
	closeStdin  bool

	awaited    bool
	foreground bool
	waitSync   sync.WaitGroup

	lock            sync.RWMutex
	gocmd           *exec.Cmd
	AutoLookPath    bool
	RaiseForBadExit bool
	mod             *LuaModule
}

func (c *Cmd) ToString() string {
	return fmt.Sprintf("%s{%s}", LuaType, strings.Join(c.gocmd.Args, ", "))
}
