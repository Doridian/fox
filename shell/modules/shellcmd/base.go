package shellcmd

import (
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

type ShellCmdModule struct {
}

func New() *ShellCmdModule {
	return &ShellCmdModule{}
}

const luaShellCmdType = "shell/modules/shellcmd/ShellCmd"

type ShellCmd struct {
	stdout *Pipe
	stderr *Pipe
	stdin  *Pipe

	gocmd            *exec.Cmd
	ErrorPropagation bool
}

func (m *ShellCmdModule) Init(L *lua.LState) {
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

	mt := L.NewTypeMetatable(luaShellCmdType)
	L.SetGlobal("cmd", mt)
	L.SetField(mt, "new", L.NewFunction(newShellCmd))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))

	funcs = map[string]lua.LGFunction{
		"read":  luaPipeRead,
		"write": luaPipeWrite,
		"close": luaPipeClose,
	}

	mt = L.NewTypeMetatable(luaShellPipeType)
	L.SetGlobal("pipe", mt)
	L.SetField(mt, "null", L.NewFunction(newNullPipe))
	L.SetField(mt, "stdin", L.NewFunction(newStdinPipe))
	L.SetField(mt, "stderr", L.NewFunction(newStderrPipe))
	L.SetField(mt, "stdout", L.NewFunction(newStdoutPipe))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))
}

func newShellCmd(L *lua.LState) int {
	return pushShellCmd(L, &ShellCmd{
		gocmd:            &exec.Cmd{},
		ErrorPropagation: false,
	})
}
