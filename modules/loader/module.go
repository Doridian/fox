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

type ModuleConfig struct {
	Global   bool
	Autoload bool
}

func DefaultLuaModuleConfig() ModuleConfig {
	return ModuleConfig{
		Global:   true,
		Autoload: true,
	}
}

type LuaModule struct {
	gomods     []*ModuleInstance
	loaderLock sync.Mutex

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

	m := &LuaModule{}

	for _, mod := range gomods {
		m.AddModuleDefault(nil, mod)
	}

	return m
}

func (m *LuaModule) AddModuleDefault(L *lua.LState, mod modules.LuaModule) {
	m.AddModule(L, mod, DefaultLuaModuleConfig())
}

func (m *LuaModule) AddModule(L *lua.LState, mod modules.LuaModule, cfg ModuleConfig) {
	m.loaderLock.Lock()
	defer m.loaderLock.Unlock()

	inst := &ModuleInstance{
		mod: mod,
		cfg: cfg,
	}
	m.gomods = append(m.gomods, inst)
	if m.loaded {
		m.preLoadMod(L, inst)
	}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	inst := &ModuleInstance{
		mod: m,
		cfg: ModuleConfig{
			Global:   true,
			Autoload: false,
		},
		loader: m.loaderInt,
	}

	return inst.loaderProxy(L)
}

func (m *LuaModule) preLoadMod(L *lua.LState, inst *ModuleInstance) {
	L.PreloadModule(inst.mod.Name(), inst.loaderProxy)

	mName := lua.LString(inst.mod.Name())
	m.builtins.Append(mName)
	if inst.cfg.Global {
		m.globals.Append(mName)
	}
	if inst.cfg.Autoload {
		m.autoload.Append(mName)
	}
}

func (m *LuaModule) loaderInt(L *lua.LState) int {
	m.loaderLock.Lock()
	defer m.loaderLock.Unlock()

	m.loaded = true

	m.builtins = L.NewTable()
	m.autoload = L.NewTable()
	m.globals = L.NewTable()

	for _, inst := range m.gomods {
		m.preLoadMod(L, inst)
	}

	for _, inst := range m.gomods {
		if inst.cfg.Autoload {
			modules.Require(L, inst.mod.Name())
		}
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
	L.PreloadModule(m.Name(), m.Loader)
	modules.Require(L, m.Name())
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt(all bool) bool {
	hit := false
	for _, inst := range m.gomods {
		if inst.mod.Interrupt(all) {
			hit = true
			if !all {
				break
			}
		}
	}
	return hit
}

func (m *LuaModule) PrePrompt() {
	for _, inst := range m.gomods {
		inst.mod.PrePrompt()
	}
}
