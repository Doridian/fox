package readline

import (
	"fmt"

	"github.com/Doridian/fox/modules/readline/config"
	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

func newReadline(L *lua.LState) int {
	prompt := L.OptString(1, "> ")
	rl, err := goreadline.New(prompt)
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

func rlSetConfig(L *lua.LState) int {
	rl, ud := Check(L, 1)
	if rl == nil {
		return 0
	}

	cfg, _ := config.Check(L, 2)
	rl.SetConfig(cfg)

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

func rlReadLineWithConfig(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}
	cfg, _ := config.Check(L, 2)

	val, err := rl.ReadLineWithConfig(cfg)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(val))
	return 1
}

func rlReadLine(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}

	val, err := rl.ReadLine()
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(val))
	return 1
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
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(val))
	return 1
}

func rlToString(L *lua.LState) int {
	rl, _ := Check(L, 1)
	if rl == nil {
		return 0
	}
	L.Push(lua.LString(fmt.Sprintf("%s{%s}", LuaType, config.ToString(rl.GetConfig()))))
	return 1
}
