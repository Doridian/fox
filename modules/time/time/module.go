package time

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.time"
const LuaTypeName = "Time"
const LuaType = LuaName + ":" + LuaTypeName

func Load(L *lua.LState, tbl *lua.LTable) {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"unixNano":  timeUnixNano,
		"unixMicro": timeUnixMicro,
		"unixMilli": timeUnixMilli,
		"unix":      timeUnix,

		"clock": timeClock,
		"date":  timeDate,

		"nanosecond": timeNanosecond,
		"second":     timeSecond,
		"minute":     timeMinute,
		"hour":       timeHour,
		"day":        timeDay,
		"month":      timeMonth,
		"year":       timeYear,

		"weekday": timeWeekday,
		"isoWeek": timeISOWeek,
		"yearDay": timeYearDay,

		"utc":    timeUTC,
		"local":  timeLocal,
		"before": timeBefore,
		"after":  timeAfter,
		"add":    timeAdd,
		"sub":    timeSub,

		"format": luaFormat,
		"string": luaString,
	}))
	mt.RawSetString("__tostring", L.NewFunction(luaToString))
	tbl.RawSetString(LuaTypeName, mt)
}
