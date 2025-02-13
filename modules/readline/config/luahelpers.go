package config

import (
	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, cfg *goreadline.Config) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = cfg
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func PushNew(L *lua.LState, cfg *goreadline.Config) int {
	if cfg == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(ToUserdata(L, cfg))
	return 1
}

func Check(L *lua.LState, i int) (*goreadline.Config, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*goreadline.Config); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return nil, nil
}
