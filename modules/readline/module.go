package readline

import (
	"github.com/Doridian/fox/modules/readline/config"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.readline"
const LuaTypeName = "ReadLine"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": newReadline,
	})

	config.Load(L, mod)

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"default": rlSetDefault,

		"config":    rlSetConfig,
		"getConfig": rlGetConfig,

		"history": rlSetHistory,

		"readLine":            rlReadLine,
		"readLineWithConfig":  rlReadLineWithConfig,
		"readLineWithDefault": rlReadLineWithDefault,

		"close": rlClose,
	}))
	mt.RawSetString("__tostring", L.NewFunction(rlToString))
	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt(all bool) bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}
