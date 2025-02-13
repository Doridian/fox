package embed

import (
	"github.com/Doridian/fox/modules/fs"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.embed"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
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
	mod.RawSetString("prefix", lua.LString(LuaName))

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

func (m *LuaModule) Interrupt(all bool) bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}
