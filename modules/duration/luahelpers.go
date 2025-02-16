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

func PushNew(L *lua.LState, d time.Duration) int {
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

func Is(L *lua.LState, i int) bool {
	ud := L.CheckUserData(i)
	_, ok := ud.Value.(time.Duration)
	return ok
}

func CheckAllowNumber(L *lua.LState, i int) (time.Duration, *lua.LUserData) {
	num, ok := L.Get(i).(lua.LNumber)
	if ok {
		return time.Duration(num), nil
	}
	return Check(L, i)
}
