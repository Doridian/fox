package duration

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.time"
const LuaTypeName = "Duration"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"nanoseconds":  durationNanoseconds,
		"milliseconds": durationMilliseconds,
		"microseconds": durationMicroseconds,
		"seconds":      durationSeconds,
		"minutes":      durationMinutes,
		"hours":        durationHours,

		"abs":    durationAbs,
		"string": luaString,
	}))
	mt.RawSetString("__tostring", L.NewFunction(luaToString))
	tbl.RawSetString(LuaTypeName, mt)
}
