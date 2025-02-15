package loader

import (
	"strings"

	"github.com/Doridian/fox/modules"
	lua "github.com/yuin/gopher-lua"
)

type ModuleInstance struct {
	mod modules.LuaModule
	cfg ModuleConfig
}

type ModuleCtor struct {
	ctor func() modules.LuaModule
	cfg  ModuleConfig
}

func (i *ModuleInstance) loaderProxy(L *lua.LState) int {
	modules.RequireDependencies(L, i.mod)
	retC := i.mod.Loader(L)
	if retC < 1 || !i.cfg.IsGlobal() {
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

	i.cfg.GlobalName = strings.TrimPrefix(strings.TrimPrefix(i.mod.Name(), "go:"), "fox.")
	return i.cfg.GlobalName
}
