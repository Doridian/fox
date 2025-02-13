package io

import (
	"io"

	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.io"
const LuaTypeName = "IO"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{})
	mod.RawSetString("SeekCurrent", lua.LNumber(io.SeekCurrent))
	mod.RawSetString("SeekStart", lua.LNumber(io.SeekStart))
	mod.RawSetString("SeekEnd", lua.LNumber(io.SeekEnd))

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), IndexFuncs()))
	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt(all bool) bool {
	return false
}
