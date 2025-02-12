package time

import (
	"github.com/Doridian/fox/modules"
	foxduration "github.com/Doridian/fox/modules/time/duration"
	foxtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.time"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	modules.RequireDependencies(L, m)

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"now": doNow,
	})
	foxduration.Load(L, mod)
	foxtime.Load(L, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Name() string {
	return LuaName
}
