package time

import (
	"github.com/Doridian/fox/modules/time/duration"
	modtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.time"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"now": doNow,
	})
	duration.Load(L, mod)
	modtime.Load(L, mod)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
