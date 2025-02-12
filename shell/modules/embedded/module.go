package embedded

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.embedded"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"loader": luaLoader,
	})
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
