package shell

import (
	"os"

	"github.com/Doridian/fox/modules/readline/config"
	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

func luaExit(L *lua.LState) int {
	exitCodeL := lua.LVAsNumber(L.CheckNumber(1))
	os.Exit(int(exitCodeL))
	return 0
}

func (s *Shell) luaGetArgs(L *lua.LState) int {
	argsL := s.l.NewTable()
	for _, arg := range s.args {
		argsL.Append(lua.LString(arg))
	}
	L.Push(argsL)
	return 1
}

func (s *Shell) luaGetArg(L *lua.LState) int {
	num := L.CheckInt(1)
	num-- // Lua is 1-based
	if num < 0 || num >= len(s.args) {
		L.Push(lua.LNil)
	} else {
		L.Push(lua.LString(s.args[num]))
	}
	return 1
}

func (s *Shell) luaPopArgs(L *lua.LState) int {
	num := L.CheckInt(1)
	if num > 0 {
		s.args = s.args[num:]
	}
	return 0
}

func (s *Shell) luaIsInteractive(L *lua.LState) int {
	L.Push(lua.LBool(s.interactive))
	return 1
}

func (s *Shell) luaSetReadlineConfig(L *lua.LState) int {
	cfg, _ := config.Check(L, 1)
	if cfg == nil {
		return 0
	}

	if s.rl == nil {
		L.Push(lua.LNil)
		return 1
	}

	cfg.Stdin = s.stdin
	cfg.Stdout = s.stdout
	cfg.Stderr = s.stderr

	rl, err := goreadline.NewFromConfig(cfg)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	s.rlLock.Lock()
	defer s.rlLock.Unlock()
	s.rl = rl
	return 0
}

func (s *Shell) luaGetReadlineConfig(L *lua.LState) int {
	if s.rl == nil {
		L.Push(lua.LNil)
		return 1
	}
	config.PushNew(L, s.rl.GetConfig())
	return 1
}
