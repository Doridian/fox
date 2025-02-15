package loader

import (
	"errors"
	"sync"

	"github.com/Doridian/fox/modules"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:loader"

type ModuleConfig struct {
	Global     *bool
	Autoload   *bool
	GlobalName string
}

type DefaultModuleConfig struct {
	Global   bool
	Autoload bool
}

func (c *ModuleConfig) IsGlobal() bool {
	if c.Global != nil {
		return *c.Global
	}
	return defaultModuleConfig.Global
}

func (c *ModuleConfig) IsAutoload() bool {
	if c.Autoload != nil {
		return *c.Autoload
	}
	return defaultModuleConfig.Autoload
}

var defaultModuleConfig = DefaultModuleConfig{
	Global:   false,
	Autoload: true,
}

func DefaultConfig() DefaultModuleConfig {
	return defaultModuleConfig
}

func SetDefaultConfig(cfg DefaultModuleConfig) {
	defaultModuleConfig = cfg
}

type LuaModule struct {
	gomods     []*ModuleInstance
	loaderLock sync.Mutex

	loaded bool

	builtins *lua.LTable
	autoload *lua.LTable
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) preLoadMod(L *lua.LState, inst *ModuleInstance) {
	L.PreloadModule(inst.mod.Name(), inst.loaderProxy)

	mName := lua.LString(inst.mod.Name())
	m.builtins.Append(mName)
	if inst.cfg.IsAutoload() {
		m.autoload.Append(mName)
	}
}

func (m *LuaModule) ManualRegisterModuleDefault(mod modules.LuaModule) error {
	return m.ManualRegisterModule(mod, ModuleConfig{})
}

func (m *LuaModule) ManualRegisterModule(mod modules.LuaModule, cfg ModuleConfig) error {
	if m.loaded {
		return errors.New("Cannot manually register modules after the loader has been loaded")
	}

	m.gomods = append(m.gomods, &ModuleInstance{
		mod: mod,
		cfg: cfg,
	})
	return nil
}

func (m *LuaModule) Loader(L *lua.LState) int {
	m.loaderLock.Lock()
	defer m.loaderLock.Unlock()

	if !m.loaded {
		ctorLock.Lock()
		for _, ctor := range ctors {
			m.gomods = append(m.gomods, &ModuleInstance{
				mod: ctor.ctor(),
				cfg: ctor.cfg,
			})
		}
		ctorLock.Unlock()

		for _, inst := range m.gomods {
			m.preLoadMod(L, inst)
		}

		for _, inst := range m.gomods {
			if inst.cfg.IsAutoload() {
				modules.Require(L, inst.mod.Name())
			}
		}
		m.loaded = true
	}

	mod := L.NewTable()
	mod.RawSetString("Builtins", m.builtins)
	mod.RawSetString("Autoload", m.autoload)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Load(L *lua.LState) {
	m.builtins = L.NewTable()
	m.autoload = L.NewTable()

	t := true
	m.preLoadMod(L, &ModuleInstance{
		mod: m,
		cfg: ModuleConfig{
			Autoload: &t,
		},
	})
	modules.Require(L, m.Name())
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt() bool {
	hit := false
	for _, inst := range m.gomods {
		if inst.mod.Interrupt() {
			hit = true
		}
	}
	return hit
}

func (m *LuaModule) PrePrompt() {
	for _, inst := range m.gomods {
		inst.mod.PrePrompt()
	}
}
