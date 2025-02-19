package cmd

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"github.com/Doridian/fox/modules/cmd/builtin"
)

// TODO: Aliasing system

const (
	ExitCodeLuaError             = -10001
	ExitCodeProcessCouldNotStart = -10002
)

type Cmd struct {
	stdout     io.Writer
	stdoutPipe io.Reader ``

	stderr     io.Writer
	stderrPipe io.Reader

	stdin     io.Reader
	stdinPipe io.Writer

	awaited    bool
	foreground bool
	waitSync   sync.WaitGroup
	startLock  sync.Mutex

	iCmd    builtin.Cmd
	iCtx    context.Context
	iCancel context.CancelFunc
	iExit   int
	iErr    error
	iDone   bool

	lock            sync.RWMutex
	gocmd           *exec.Cmd
	AutoLookPath    bool
	RaiseForBadExit bool
	mod             *LuaModule

	closeQueue []io.Closer
}

func (c *Cmd) ToString() string {
	return fmt.Sprintf("%s{%s}", LuaType, strings.Join(c.gocmd.Args, ", "))
}
