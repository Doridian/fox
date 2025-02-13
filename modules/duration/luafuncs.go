package duration

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func durationHours(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LNumber(d.Hours()))
	return 1
}

func durationMinutes(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LNumber(d.Minutes()))
	return 1
}

func durationSeconds(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LNumber(d.Seconds()))
	return 1
}

func durationMilliseconds(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LNumber(d.Milliseconds()))
	return 1
}

func durationMicroseconds(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LNumber(d.Microseconds()))
	return 1
}

func durationNanoseconds(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LNumber(d.Nanoseconds()))
	return 1
}

func durationAbs(L *lua.LState) int {
	d, _ := Check(L, 1)
	return Push(L, d.Abs())
}

func luaEq(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := Check(L, 2)
	L.Push(lua.LBool(d == d2))
	return 1
}

func luaLt(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := Check(L, 2)
	L.Push(lua.LBool(d > d2))
	return 1
}

func luaLe(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := Check(L, 2)
	L.Push(lua.LBool(d <= d2))
	return 1
}

func luaAdd(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := CheckAllowNumber(L, 2)
	return Push(L, (d + d2))
}

func luaSub(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := CheckAllowNumber(L, 2)
	return Push(L, (d - d2))
}

func luaMul(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := CheckAllowNumber(L, 2)
	return Push(L, (d * d2))
}

func luaDiv(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := CheckAllowNumber(L, 2)
	return Push(L, (d / d2))
}

func luaUnm(L *lua.LState) int {
	d, _ := Check(L, 1)
	return Push(L, -d)
}

func luaString(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LString(d.String()))
	return 1
}

func luaToString(L *lua.LState) int {
	d, _ := Check(L, 1)
	L.Push(lua.LString(fmt.Sprintf("%s{%s}", LuaType, d.String())))
	return 1
}
