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
	gName := i.globalName()

	gLastDot := strings.LastIndex(gName, ".")
	if gLastDot <= 0 {
		L.SetGlobal(gName, modL)
		return retC
	}

	gTbl := gName[:gLastDot]
	gName = gName[gLastDot+1:]

	tbl := L.FindTable(L.G.Global, gTbl, 1).(*lua.LTable)
	tbl.RawSetString(gName, modL)

	return retC
}

func (i *ModuleInstance) globalName() string {
	if i.cfg.GlobalName != "" {
		return i.cfg.GlobalName
	}

	return strings.TrimPrefix(i.mod.Name(), "fox.")
}
