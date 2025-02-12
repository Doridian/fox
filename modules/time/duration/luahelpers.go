package duration

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, d time.Duration) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = d
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func Push(L *lua.LState, d time.Duration) int {
	L.Push(ToUserdata(L, d))
	return 1
}

func Check(L *lua.LState, i int) (time.Duration, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(time.Duration); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return 0, nil
}
