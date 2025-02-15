package duration

import (
	"time"

	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:duration"
const LuaTypeName = "Duration"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func newLuaModule() modules.LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"parse": durationParse,
	})

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"nanoseconds":  durationNanoseconds,
		"milliseconds": durationMilliseconds,
		"microseconds": durationMicroseconds,
		"seconds":      durationSeconds,
		"minutes":      durationMinutes,
		"hours":        durationHours,

		"sleepFor": durationSleepFor,

		"abs":    durationAbs,
		"string": luaString,
	}))

	mt.RawSetString("Nanosecond", ToUserdata(L, time.Nanosecond))
	mt.RawSetString("Microsecond", ToUserdata(L, time.Microsecond))
	mt.RawSetString("Millisecond", ToUserdata(L, time.Millisecond))
	mt.RawSetString("Second", ToUserdata(L, time.Second))
	mt.RawSetString("Minute", ToUserdata(L, time.Minute))
	mt.RawSetString("Hour", ToUserdata(L, time.Hour))

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

	mod.RawSetString(LuaTypeName, mt)
	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{}
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
