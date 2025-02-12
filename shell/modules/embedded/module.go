package embedded

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.embedded"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.NewTable()
	L.SetFuncs(mod, map[string]lua.LGFunction{
		"bareLoader":   luaBareLoader,
		"prefixLoader": luaPrefixLoader,
		"readFile":     luaReadFile,
		"doFile":       luaDoFile,
		"loadFile":     luaLoadFile,
	}, mod)
	mod.RawSetString("path", lua.LString("root/?.lua;root/?/init.lua"))
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
