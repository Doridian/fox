package embed

import (
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

		"readDir": luaReadDir,
	}, mod)
	mod.RawSetString("path", lua.LString("root/?.lua;root/?/init.lua"))
	mod.RawSetString("prefix", lua.LString(LuaName))

	packagesL := L.GetGlobal("package").(*lua.LTable)
	loadersL := packagesL.RawGetString("loaders").(*lua.LTable)
	loadersL.Append(mod.RawGetString("loader"))

	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
