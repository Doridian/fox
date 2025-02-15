package embed

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/fs"
	"github.com/Doridian/fox/modules/loader"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:embed"

type LuaModule struct {
}

func newLuaModule() modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.NewTable()
	L.SetFuncs(mod, map[string]lua.LGFunction{
		"loader": luaLoader,

		"readFile": luaReadFile,
		"doFile":   luaDoFile,
		"loadFile": luaLoadFile,
		"openFile": luaOpenFile,

		"readDir": luaReadDir,
	}, mod)
	mod.RawSetString("path", lua.LString("root/?.lua;root/?/init.lua"))
	mod.RawSetString("prefix", lua.LString("embed:"))

	loadersL := L.GetField(L.GetField(L.Get(lua.EnvironIndex), "package"), "loaders")
	if loadersL == nil || loadersL == lua.LNil {
		L.Push(mod)
		return 1
	}
	loadersL.(*lua.LTable).Append(mod.RawGetString("loader"))

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{fs.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt() bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}

func init() {
	loader.AddModuleDefault(newLuaModule)
}
