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
