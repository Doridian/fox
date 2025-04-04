package readline

import (
	"fmt"

	"github.com/Doridian/fox/modules/readline/config"
	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

func (m *LuaModule) newReadline(L *lua.LState) int {
	cfg := &goreadline.Config{
		Prompt: L.OptString(1, "> "),
	}
	cfg = fixConfig(m.loader, cfg)

	rl, err := goreadline.NewFromConfig(cfg)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return PushNew(L, rl)
}

func (m *LuaModule) newReadlineFromConfig(L *lua.LState) int {
	cfg, _ := config.Check(L, 1)
	if cfg == nil {
		return 0
	}

	cfg = fixConfig(m.loader, cfg)

	rl, err := goreadline.NewFromConfig(cfg)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return PushNew(L, rl)
}

func rlSetDefault(L *lua.LState) int {
	rl, ud := Check(L, 1)
	if rl == nil {
		return 0
	}

	val := L.CheckString(2)
	rl.SetDefault(val)

	L.Push(ud)
	return 1
}

func rlSetHistory(L *lua.LState) int {
	rl, ud := Check(L, 1)
	if rl == nil {
		return 0
	}

	val := L.CheckBool(2)
	if val {
		rl.EnableHistory()
	} else {
		rl.DisableHistory()
	}

	L.Push(ud)
	return 1
}

func (m *LuaModule) rlSetConfig(L *lua.LState) int {
	rl, ud := Check(L, 1)
	if rl == nil {
		return 0
	}

	cfg, _ := config.Check(L, 2)
	cfg = fixConfig(m.loader, cfg)

	err := rl.SetConfig(cfg)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(ud)
	return 1
}

func rlGetConfig(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}

	return config.PushNew(L, rl.GetConfig())
}

func (m *LuaModule) rlReadLineWithConfig(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}
	cfg, _ := config.Check(L, 2)
	cfg = fixConfig(m.loader, cfg)

	val, err := rl.ReadLineWithConfig(cfg)
	return rlResultHandler(L, val, err)
}

func rlReadLine(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}

	val, err := rl.ReadLine()
	return rlResultHandler(L, val, err)
}

func rlReadLineWithDefault(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}
	def := L.CheckString(2)
	if def == "" {
		L.ArgError(2, "default value must not be empty")
		return 0
	}

	val, err := rl.ReadLineWithDefault(def)
	return rlResultHandler(L, val, err)
}

func rlToString(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}
	L.Push(lua.LString(fmt.Sprintf("%s{%s}", LuaType, config.ToString(rl.GetConfig()))))
	return 1
}

func rlClose(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}

	_ = rl.Close()
	return 0
}
