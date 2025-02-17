package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/Doridian/fox/modules/cmd/integrated"
)

const (
	ExitCodeLuaError             = -10001
	ExitCodeProcessCouldNotStart = -10002
)

type Cmd struct {
	stdout      interface{}
	closeStdout bool
	stderr      interface{}
	closeStderr bool
	stdin       interface{}
	closeStdin  bool

	awaited    bool
	foreground bool
	waitSync   sync.WaitGroup
	startLock  sync.Mutex

	iCmd    integrated.Cmd
	iCtx    context.Context
	iCancel context.CancelFunc
	iExit   int

	lock            sync.RWMutex
	gocmd           *exec.Cmd
	AutoLookPath    bool
	RaiseForBadExit bool
	mod             *LuaModule
}

func (c *Cmd) ToString() string {
	return fmt.Sprintf("%s{%s}", LuaType, strings.Join(c.gocmd.Args, ", "))
}
