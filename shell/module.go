package shell

import (
	"context"
	"os"
	"sync"

	"github.com/Doridian/fox/modules"
	"github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.shell"

type Shell struct {
	args []string

	l     *lua.LState
	lLock sync.Mutex

	mod   *lua.LTable
	print *lua.LFunction

	ctx       context.Context
	cancelCtx context.CancelFunc

	mainMod modules.LuaModule

	rlLock sync.Mutex
	rl     *readline.Instance

	signals chan os.Signal
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
