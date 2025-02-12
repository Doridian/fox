package modules

import lua "github.com/yuin/gopher-lua"

type LuaModule interface {
	Name() string
	Loader(l *lua.LState) int
}

func Preload(mod LuaModule, L *lua.LState) {
	L.PreloadModule(mod.Name(), mod.Loader)
}
