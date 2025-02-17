package cmd

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	lua "github.com/yuin/gopher-lua"
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

	lock            sync.RWMutex
	gocmd           *exec.Cmd
	AutoLookPath    bool
	RaiseForBadExit bool
	mod             *LuaModule
	preReqs         []*lua.LFunction
}

func (c *Cmd) ToString() string {
	return fmt.Sprintf("%s{%s}", LuaType, strings.Join(c.gocmd.Args, ", "))
}
