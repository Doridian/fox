package readline

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.readline"
const LuaTypeName = "ReadLine"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": newReadline,
	})

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"default": rlSetDefault,

		"config":    rlSetConfig,
		"getConfig": rlGetConfig,
	}))
	mt.RawSetString("__tostring", L.NewFunction(rlToString))
	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}
