package shellcmd

import (
	"os"
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
	stdout *ShellPipe
	stderr *ShellPipe
	stdin  *ShellPipe

	gocmd            *exec.Cmd
	ErrorPropagation bool
}

func (m *ShellCmdModule) Init(L *lua.LState) {
	funcs := map[string]lua.LGFunction{
		"path": getSetPath,
		"dir":  getSetDir,
		"args": getSetArgs,
		"env":  getSetEnv,

		"stdout": getSetStdout,
		"stderr": getSetStderr,
		"stdin":  getSetStdin,

		"run":   doRun,
		"start": doStart,
		"wait":  doWait,

		"errorPropagation": getSetErrorPropagation,
	}

	mt := L.NewTypeMetatable(luaShellCmdType)
	L.SetGlobal("shellcmd", mt)
	L.SetField(mt, "new", L.NewFunction(newShellCmd))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))

	funcs = map[string]lua.LGFunction{}

	mt = L.NewTypeMetatable(luaShellPipeType)
	L.SetGlobal("shellpipe", mt)
	L.SetField(mt, "null", L.NewFunction(newNullShellPipe))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))
}

func newShellCmd(L *lua.LState) int {
	return pushShellCmd(L, &ShellCmd{
		gocmd: &exec.Cmd{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			Stdin:  os.Stdin,
		},
		ErrorPropagation: false,
	})
}
