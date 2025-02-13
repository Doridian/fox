package env

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.Env"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
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

func (m *LuaModule) Interrupt(all bool) bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}
