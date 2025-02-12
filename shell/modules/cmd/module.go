package cmd

import (
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "fox.cmd"
const LuaType = LuaName + ":Cmd"

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

	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new": newCmd,
	})
	L.Push(mod)
	return 1
}

func (m *LuaModule) Name() string {
	return LuaName
}
