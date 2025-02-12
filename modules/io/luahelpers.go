package io

import (
	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, f interface{}) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = f
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func Push(L *lua.LState, f interface{}) int {
	if f == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(ToUserdata(L, f))
	return 1
}

func Check(L *lua.LState, i int) (interface{}, *lua.LUserData) {
	ud := L.CheckUserData(i)
	return ud.Value, ud
}

func IndexFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"close": ioClose,
		"read":  ioRead,
		"write": ioWrite,
		"seek":  ioSeek,
	}
}
