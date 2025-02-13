package loader

import (
	"github.com/Doridian/fox/modules"
	lua "github.com/yuin/gopher-lua"
)

type moduleProxyInt struct {
	base modules.LuaModule

	autoload bool
	global   bool
}

func proxyGoMod(base modules.LuaModule) *moduleProxyInt {
	baseProxyInt, ok := base.(*moduleProxyInt)
	if ok {
		return baseProxyInt
	}

	global := true
	autoload := true
	baseProxy, ok := base.(ModuleProxy)
	if ok {
		global = baseProxy.Global()
		autoload = baseProxy.Autoload()
		baseProxyWithBase, ok := base.(ModuleProxyWithBase)
		if ok {
			base = baseProxyWithBase.Base()
		}
	}

	return &moduleProxyInt{
		base:     base,
		autoload: autoload,
		global:   global,
	}
}

func (m *moduleProxyInt) Loader(L *lua.LState) int {
	return loaderViaProxy(L, m, m.base.Loader)
}

func (m *moduleProxyInt) Dependencies() []string {
	return m.base.Dependencies()
}

func (m *moduleProxyInt) Name() string {
	return m.base.Name()
}

func (m *moduleProxyInt) Global() bool {
	return m.global
}

func (m *moduleProxyInt) SetGlobal(global bool) {
	m.global = global
}

func (m *moduleProxyInt) Autoload() bool {
	return m.autoload
}

func (m *moduleProxyInt) SetAutoload(autoload bool) {
	m.autoload = autoload
}

func (m *moduleProxyInt) Base() modules.LuaModule {
	return m.base
}

func (m *moduleProxyInt) Interrupt(all bool) bool {
	return m.base.Interrupt(all)
}
