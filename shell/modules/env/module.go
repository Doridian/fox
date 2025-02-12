package env

import (
	lua "github.com/yuin/gopher-lua"
)

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"__index":    envIndex,
		"__newindex": envNewIndex,
	})
	L.SetMetatable(mod, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return "env"
}
