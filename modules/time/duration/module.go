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

	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__add": luaAdd,
		"__sub": luaSub,
		"__mul": luaMul,
		"__div": luaDiv,
		"__unm": luaUnm,

		"__eq": luaEq,
		"__lt": luaLt,
		"__le": luaLe,

		"__tostring": luaToString,
	})

	tbl.RawSetString(LuaTypeName, mt)
}
