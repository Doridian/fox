package loader

import (
	"sync"

	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/cmd"
	"github.com/Doridian/fox/modules/duration"
	"github.com/Doridian/fox/modules/embed"
	"github.com/Doridian/fox/modules/env"
	"github.com/Doridian/fox/modules/fs"
	"github.com/Doridian/fox/modules/io"
	"github.com/Doridian/fox/modules/pipe"
	"github.com/Doridian/fox/modules/readline"
	"github.com/Doridian/fox/modules/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.index"

type LuaModule struct {
	gomods     []*moduleProxyInt
	loaderLock sync.Mutex

	global bool
	loaded bool

	builtins *lua.LTable
	autoload *lua.LTable
	globals  *lua.LTable
}

func NewLuaModule() *LuaModule {
	gomods := []modules.LuaModule{
		time.NewLuaModule(),
		duration.NewLuaModule(),
		io.NewLuaModule(),
		fs.NewLuaModule(),
		embed.NewLuaModule(),
		env.NewLuaModule(),
		pipe.NewLuaModule(),
		cmd.NewLuaModule(),
		readline.NewLuaModule(),
	}

	gomodsProxied := make([]*moduleProxyInt, 0, len(gomods))
	for _, m := range gomods {
		gomodsProxied = append(gomodsProxied, proxyGoMod(m))
	}

	return &LuaModule{
		gomods: gomodsProxied,
		global: true,
	}
}

func (m *LuaModule) AddModule(L *lua.LState, mod modules.LuaModule) {
	m.loaderLock.Lock()
	defer m.loaderLock.Unlock()

	pm := proxyGoMod(mod)
	m.gomods = append(m.gomods, pm)
	if m.loaded {
		pm.SetAutoload(false)
		m.preLoadMod(L, pm)
		m.postLoadMod(L, pm)
	}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	return loaderViaProxy(L, m, m.loaderInt)
}

func (m *LuaModule) preLoadMod(L *lua.LState, mod ModuleProxy) {
	modules.Preload(L, mod)

	mName := lua.LString(m.Name())
	m.builtins.Append(mName)
	if mod.Global() {
		m.globals.Append(mName)
	}
	if mod.Autoload() {
		m.autoload.Append(mName)
	}
}

func (m *LuaModule) postLoadMod(L *lua.LState, mod ModuleProxy) {
	if mod.Autoload() {
		modules.Require(L, mod.Name())
	}
}

func (m *LuaModule) loaderInt(L *lua.LState) int {
	m.loaderLock.Lock()
	defer m.loaderLock.Unlock()

	m.loaded = true

	m.builtins = L.NewTable()
	m.autoload = L.NewTable()
	m.globals = L.NewTable()

	for _, mod := range m.gomods {
		m.preLoadMod(L, mod)
	}

	for _, mod := range m.gomods {
		m.postLoadMod(L, mod)
	}

	mod := L.NewTable()
	mod.RawSetString("BuiltIns", m.builtins)
	mod.RawSetString("AutoLoad", m.autoload)
	mod.RawSetString("Globals", m.globals)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Load(L *lua.LState) {
	modules.Preload(L, m)
	modules.Require(L, m.Name())
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Global() bool {
	return m.global
}

func (m *LuaModule) SetGlobal(global bool) {
	m.global = global
}

func (m *LuaModule) Autoload() bool {
	return true
}

func (m *LuaModule) SetAutoload(autoload bool) {
	if !autoload {
		panic("cannot disable autoload for the loader module")
	}
}

func (m *LuaModule) Interrupt(all bool) bool {
	hit := false
	for _, m := range m.gomods {
		if m.Interrupt(all) {
			hit = true
			if !all {
				break
			}
		}
	}
	return hit
}

func (m *LuaModule) PrePrompt() {
	for _, m := range m.gomods {
		m.PrePrompt()
	}
}
