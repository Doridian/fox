package cmd

import (
	"context"
	"fmt"
	"io"
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
	stdout       interface{}
	stdoutCloser io.Closer
	stderr       interface{}
	stderrCloser io.Closer
	stdin        interface{}
	stdinCloser  io.Closer

	awaited    bool
	foreground bool
	waitSync   sync.WaitGroup
	startLock  sync.Mutex

	iCmd     integrated.Cmd
	iCtx     context.Context
	iCancel  context.CancelFunc
	iExit    int
	iErr     error
	iCmdWait sync.WaitGroup

	lock            sync.RWMutex
	gocmd           *exec.Cmd
	AutoLookPath    bool
	RaiseForBadExit bool
	mod             *LuaModule
}

func (c *Cmd) ToString() string {
	return fmt.Sprintf("%s{%s}", LuaType, strings.Join(c.gocmd.Args, ", "))
}
