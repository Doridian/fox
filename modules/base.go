package modules

import (
	lua "github.com/yuin/gopher-lua"
)

type LuaModule interface {
	Name() string
	Dependencies() []string
	Loader(l *lua.LState) int
	Interrupt(all bool) bool // Return true if you took the Interrupt (stops trying to stop other things unless all is true)
	PrePrompt()
}

func Preload(L *lua.LState, m LuaModule) {
	L.PreloadModule(m.Name(), m.Loader)
}

func RequireDependencies(L *lua.LState, m LuaModule) {
	for _, dep := range m.Dependencies() {
		Require(L, dep)
	}
}

func Require(L *lua.LState, v string) lua.LValue {
	requireL := L.GetGlobal("require")
	L.Push(requireL)
	L.Push(lua.LString(v))
	L.Call(1, 1)
	mod := L.Get(-1)
	L.Pop(1)
	return mod
}
