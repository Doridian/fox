package embedded

import (
	lua "github.com/yuin/gopher-lua"
)

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Init(L *lua.LState) {
	loader := L.NewFunction(luaLoader)

	packageL := L.GetGlobal("package").(*lua.LTable)
	loadersL := L.GetField(packageL, "preload").(*lua.LTable)

	loadersL.Append(loader)
}
