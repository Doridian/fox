package modules

import (
	"log"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

type LuaModule interface {
	Name() string
	Dependencies() []string
	Loader(l *lua.LState) int
}

func Preload(L *lua.LState, m LuaModule) {
	L.PreloadModule(m.Name(), m.Loader)
}

func RequireDependencies(L *lua.LState, m LuaModule) {
	log.Printf("RequireDependencies(%s)", m.Name())
	for _, dep := range m.Dependencies() {
		Require(L, dep)
	}
}

func Require(L *lua.LState, v string) lua.LValue {
	log.Printf("Require(%s)", v)
	requireL := L.GetGlobal("require")
	L.Push(requireL)
	L.Push(lua.LString(v))
	L.Call(1, 1)
	mod := L.Get(-1)
	L.Pop(1)
	return mod
}

func RequireGlobal(L *lua.LState, m LuaModule) lua.LValue {
	mod := Require(L, m.Name())

	modVarName := m.Name()
	modDotIdx := strings.LastIndex(modVarName, ".")
	if modDotIdx != -1 {
		modVarName = modVarName[modDotIdx+1:]
	}
	L.SetGlobal(modVarName, mod)

	return mod
}
