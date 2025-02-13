package duration

import (
	"time"

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

	mt.RawSetString("nanosecond", ToUserdata(L, time.Nanosecond))
	mt.RawSetString("microsecond", ToUserdata(L, time.Microsecond))
	mt.RawSetString("millisecond", ToUserdata(L, time.Millisecond))
	mt.RawSetString("second", ToUserdata(L, time.Second))
	mt.RawSetString("minute", ToUserdata(L, time.Minute))
	mt.RawSetString("hour", ToUserdata(L, time.Hour))

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

		"parse": durationParse,
	})

	tbl.RawSetString(LuaTypeName, mt)
}
