package cmd

import (
	lua "github.com/yuin/gopher-lua"
)

const luaCmdType = "go://fox/cmd/Cmd"

type LuaModule struct {
}

func NewLuaModule() *LuaModule {
	return &LuaModule{}
}

func (m *LuaModule) Init(L *lua.LState) {
	funcs := map[string]lua.LGFunction{
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
	}

	mt := L.NewTypeMetatable(luaCmdType)
	L.SetGlobal("cmd", mt)
	L.SetField(mt, "new", L.NewFunction(newCmd))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))
}
