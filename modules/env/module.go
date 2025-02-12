package env

import (
	"github.com/Doridian/fox/modules"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.env"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	modules.RequireDependencies(L, m)

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"__index":    envIndex,
		"__newindex": envNewIndex,
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
