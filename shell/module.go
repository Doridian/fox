package shell

import (
	"context"
	"io"
	"sync"

	"github.com/Doridian/fox/modules"
	"github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:shell"

type Shell struct {
	args        []string
	interactive bool

	l     *lua.LState
	lLock sync.Mutex

	mod *lua.LTable

	ctx       context.Context
	cancelCtx context.CancelFunc

	mainMod modules.LuaModule

	rlLock sync.Mutex
	rl     *readline.Instance

	ShowErrors bool

	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func (s *Shell) Dependencies() []string {
	return []string{}
}

func (s *Shell) Interrupt() bool {
	return false
}

func (s *Shell) Name() string {
	return LuaName
}

func (s *Shell) PrePrompt() {
	// no-op
}
