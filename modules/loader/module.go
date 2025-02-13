package loader

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/cmd"
	"github.com/Doridian/fox/modules/duration"
	"github.com/Doridian/fox/modules/embed"
	"github.com/Doridian/fox/modules/env"
	"github.com/Doridian/fox/modules/fs"
	"github.com/Doridian/fox/modules/io"
	"github.com/Doridian/fox/modules/pipe"
	"github.com/Doridian/fox/modules/time"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.index"

type LuaModule struct {
	gomods []*moduleProxyInt

	global bool
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

func (m *LuaModule) Loader(L *lua.LState) int {
	return loaderViaProxy(L, m, m.loaderInt)
}

func (m *LuaModule) loaderInt(L *lua.LState) int {
	builtins := L.NewTable()
	autoload := L.NewTable()
	globals := L.NewTable()

	for _, m := range m.gomods {
		modules.Preload(L, m)

		mName := lua.LString(m.Name())
		builtins.Append(mName)
		if m.global {
			globals.Append(mName)
		}
		if m.autoload {
			autoload.Append(mName)
		}
	}

	for _, m := range m.gomods {
		if !m.autoload {
			continue
		}
		modules.Require(L, m.Name())
	}

	mod := L.NewTable()
	mod.RawSetString("builtins", builtins)
	mod.RawSetString("autoload", autoload)
	mod.RawSetString("globals", globals)
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
