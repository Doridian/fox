package vars

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:vars"

var varTable = make(map[string]lua.LString)

type LuaModule struct {
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"__index":    varsIndex,
		"__newindex": varsNewIndex,
		"__call":     varsCall,
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
	loader.AddModule(newLuaModule, loader.ModuleConfig{
		GlobalName: "V",
	})
}
