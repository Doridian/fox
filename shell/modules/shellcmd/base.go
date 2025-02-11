package shellcmd

import (
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

type ShellCmdModule struct {
}

const luaShellCmdType = "shell/modules/shellcmd"

type ShellCmd struct {
	Stdin            *ShellCmd
	Stdout           *ShellCmd
	Stderr           *ShellCmd
	Gocmd            *exec.Cmd
	ErrorPropagation bool
}

func New() *ShellCmdModule {
	return &ShellCmdModule{}
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
}

func newShellCmd(L *lua.LState) int {
	return pushShellCmd(L, &ShellCmd{
		Gocmd:            &exec.Cmd{},
		ErrorPropagation: false,
	})
}
