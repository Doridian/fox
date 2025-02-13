package readline

import (
	goreadline "github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, rl *goreadline.Instance) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = rl
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func PushNew(L *lua.LState, rl *goreadline.Instance) int {
	if rl == nil {
		L.Push(lua.LNil)
		return 1
	}
	L.Push(ToUserdata(L, rl))
	return 1
}

func Check(L *lua.LState, i int) (*goreadline.Instance, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*goreadline.Instance); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return nil, nil
}
