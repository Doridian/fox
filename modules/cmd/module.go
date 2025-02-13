package cmd

import (
	"github.com/Doridian/fox/modules/pipe"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.cmd"
const LuaTypeName = "Cmd"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"dir": getSetDir,
		"cmd": getSetCmd,
		"env": getSetEnv,

		"stdout":     getSetStdout,
		"stdoutPipe": getStdoutPipe,
		"stderr":     getSetStderr,
		"stderrPipe": getStderrPipe,
		"stdin":      getSetStdin,
		"stdinPipe":  getStdinPipe,

		"run":   doRun,
		"start": doStart,
		"wait":  doWait,

		"errorPropagation": getSetErrorPropagation,
	}))
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__tostring": cmdToString,
		"__call":     doRun,
	})

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": newCmd,
	})

	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{pipe.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}
