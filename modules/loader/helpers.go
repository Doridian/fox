package loader

import (
	"strings"

	"github.com/Doridian/fox/modules"
	lua "github.com/yuin/gopher-lua"
)

type ModuleProxy interface {
	modules.LuaModule

	Global() bool
	SetGlobal(bool)
	Autoload() bool
	SetAutoload(bool)
}

type ModuleProxyWithBase interface {
	modules.LuaModule

	Global() bool
	SetGlobal(bool)
	Base() modules.LuaModule
}

func loaderViaProxy(L *lua.LState, m ModuleProxy, baseLoader lua.LGFunction) int {
	modules.RequireDependencies(L, m)
	retC := baseLoader(L)
	if retC < 1 || !m.Global() {
		return retC
	}

	modL := L.Get(-1)

	modVarName := m.Name()
	modDotIdx := strings.LastIndex(modVarName, ".")
	if modDotIdx != -1 {
		modVarName = modVarName[modDotIdx+1:]
	}
	L.SetGlobal(modVarName, modL)

	return retC
}
