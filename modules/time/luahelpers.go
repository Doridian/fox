package time

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

func ToUserdata(L *lua.LState, t time.Time) *lua.LUserData {
	ud := L.NewUserData()
	ud.Value = t
	L.SetMetatable(ud, L.GetTypeMetatable(LuaType))
	return ud
}

func PushNew(L *lua.LState, t time.Time) int {
	L.Push(ToUserdata(L, t))
	return 1
}

func Is(L *lua.LState, i int) bool {
	ud := L.CheckUserData(i)
	_, ok := ud.Value.(time.Time)
	return ok
}

func Check(L *lua.LState, i int) (time.Time, *lua.LUserData) {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(time.Time); ok {
		return v, ud
	}

	L.ArgError(i, LuaType+" expected")
	return time.Time{}, nil
}
