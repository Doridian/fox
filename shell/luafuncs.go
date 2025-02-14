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

func (s *Shell) luaSetReadlineConfig(L *lua.LState) int {
	cfg, _ := config.Check(L, 1)
	if cfg == nil {
		return 0
	}

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
	config.PushNew(L, s.rl.GetConfig())
	return 1
}
