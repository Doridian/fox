package env

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:env"

type LuaModule struct {
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"__index":    envIndex,
		"__newindex": envNewIndex,
		"__call":     envCall,
	})
	L.SetMetatable(mod, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
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
