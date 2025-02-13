package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/Doridian/fox/modules/pipe"
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
	AutoLookPath     bool
	ErrorPropagation bool
}

func (c *Cmd) ToString() string {
	return fmt.Sprintf("%s{%s}", LuaType, strings.Join(c.gocmd.Args, ", "))
}
