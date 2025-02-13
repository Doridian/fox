package shell

import (
	"os"

	"github.com/Doridian/fox/modules/readline/config"
	lua "github.com/yuin/gopher-lua"
)

func luaExit(L *lua.LState) int {
	exitCodeL := lua.LVAsNumber(L.CheckNumber(1))
	os.Exit(int(exitCodeL))
	return 0
}

func (s *Shell) luaSetReadLineConfig(L *lua.LState) int {
	cfg, _ := config.Check(L, 1)
	if cfg == nil {
		return 0
	}

	s.rl.SetConfig(cfg)
	return 0
}

func (s *Shell) luaGetReadLineConfig(L *lua.LState) int {
	config.PushNew(L, s.rl.GetConfig())
	return 1
}
