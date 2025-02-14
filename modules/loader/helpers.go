package loader

import (
	"strings"

	"github.com/Doridian/fox/modules"
	lua "github.com/yuin/gopher-lua"
)

type ModuleInstance struct {
	mod modules.LuaModule
	cfg ModuleConfig

	loader lua.LGFunction
}

func (i *ModuleInstance) loaderProxy(L *lua.LState) int {
	modules.RequireDependencies(L, i.mod)
	var retC int
	if i.loader != nil {
		retC = i.loader(L)
	} else {
		retC = i.mod.Loader(L)
	}
	if retC < 1 || !i.cfg.Global {
		return retC
	}

	modL := L.Get(-1)

	modVarName := i.mod.Name()
	modDotIdx := strings.LastIndex(modVarName, ".")
	if modDotIdx != -1 {
		modVarName = modVarName[modDotIdx+1:]
	}
	L.SetGlobal(modVarName, modL)

	return retC
}
