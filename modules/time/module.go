package time

import (
	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/duration"
	"github.com/Doridian/fox/modules/loader"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:time"
const LuaTypeName = "Time"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func newLuaModule() modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"now":   timeNow,
		"parse": timeParse,
	})

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"unixNano":  timeUnixNano,
		"unixMicro": timeUnixMicro,
		"unixMilli": timeUnixMilli,
		"unix":      timeUnix,

		"clock": timeClock,
		"date":  timeDate,

		"addDate": timeAddDate,

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

		"utc":   timeUTC,
		"local": timeLocal,
		"delta": timeDelta,
		"until": timeUntil,
		"since": timeSince,

		"sleepUntil": timeSleepUntil,

		"format": luaFormat,
		"string": luaString,
	}))

	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__add": timeAddDuration,
		"__sub": timeSubDuration,

		"__eq": timeEq,
		"__lt": timeBefore,
		"__le": timeNotAfter,

		"__tostring": luaToString,
	})

	mod.RawSetString(LuaTypeName, mt)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{duration.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt() bool {
	return false
}

func (m *LuaModule) PrePrompt() {
	// no-op
}

func init() {
	loader.AddModuleDefault(newLuaModule)
}
