package readline

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/modules/readline/config"
	"github.com/Doridian/fox/shell"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:readline"
const LuaTypeName = "Readline"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
	loader *loader.LuaModule
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{
		loader: loader,
	}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new":           m.newReadline,
		"newFromConfig": m.newReadlineFromConfig,
	})

	config.Load(L, mod)

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"default": rlSetDefault,

		"config":    m.rlSetConfig,
		"getConfig": rlGetConfig,

		"history": rlSetHistory,

		"readLine":            rlReadLine,
		"readLineWithConfig":  m.rlReadLineWithConfig,
		"readLineWithDefault": rlReadLineWithDefault,

		"close": rlClose,
	}))
	mt.RawSetString("__tostring", L.NewFunction(rlToString))
	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{shell.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt() bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}

func init() {
	loader.AddModuleDefault(newLuaModule)
}
