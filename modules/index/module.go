package index

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/cmd"
	foxembed "github.com/Doridian/fox/modules/embed"
	"github.com/Doridian/fox/modules/env"
	foxfs "github.com/Doridian/fox/modules/fs"
	foxio "github.com/Doridian/fox/modules/io"
	"github.com/Doridian/fox/modules/pipe"
	foxtime "github.com/Doridian/fox/modules/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.index"

type LuaModule struct {
	gomods map[modules.LuaModule]bool
}

func NewLuaModule() *LuaModule {
	return &LuaModule{
		gomods: map[modules.LuaModule]bool{
			foxtime.NewLuaModule():  true,
			foxio.NewLuaModule():    true,
			foxfs.NewLuaModule():    true,
			foxembed.NewLuaModule(): true,
			env.NewLuaModule():      true,
			pipe.NewLuaModule():     true,
			cmd.NewLuaModule():      true,
		},
	}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	modules.RequireDependencies(L, m)

	builtins := L.NewTable()
	defaults := L.NewTable()

	for m := range m.gomods {
		modules.Preload(L, m)
		builtins.Append(lua.LString(m.Name()))
	}

	for m, autoload := range m.gomods {
		if !autoload {
			continue
		}
		modules.RequireGlobal(L, m)
		defaults.Append(lua.LString(m.Name()))
	}

	mod := L.NewTable()
	mod.RawSetString("builtins", builtins)
	mod.RawSetString("defaults", defaults)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Load(L *lua.LState) {
	modules.Preload(L, m)
	modules.RequireGlobal(L, m)
}

func (m *LuaModule) Name() string {
	return LuaName
}
