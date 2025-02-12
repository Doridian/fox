package env

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaType = "go://fox/env/Env"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Init(L *lua.LState) {
	mt := L.NewTypeMetatable(LuaType)
	L.SetGlobal("env", mt)
	// TODO: __pairs
	L.SetField(mt, "__index", L.NewFunction(envIndex))
	L.SetField(mt, "__newindex", L.NewFunction(envNewIndex))
	L.SetMetatable(mt, mt)
}
