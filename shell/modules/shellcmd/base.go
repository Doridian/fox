package shellcmd

import (
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

type ShellCmdModule struct {
}

const luaShellCmdType = "shell/modules/shellcmd"

type ShellCmd struct {
	RedirectStdin  *ShellCmd
	RedirectStdout *ShellCmd
	RedirectStderr *ShellCmd
	Gocmd          *exec.Cmd
}

func New() *ShellCmdModule {
	return &ShellCmdModule{}
}

func (m *ShellCmdModule) Init(L *lua.LState) {
	funcs := map[string]lua.LGFunction{
		"path": getSetPath,
		"dir":  getSetDir,
		"args": getSetArgs,
	}

	mt := L.NewTypeMetatable(luaShellCmdType)
	L.SetGlobal("shellcmd", mt)
	L.SetField(mt, "new", L.NewFunction(newShellCmd))
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), funcs))
}

func newShellCmd(L *lua.LState) int {
	cmd := &ShellCmd{
		Gocmd: &exec.Cmd{},
	}
	ud := L.NewUserData()
	ud.Value = cmd
	L.SetMetatable(ud, L.GetTypeMetatable(luaShellCmdType))
	L.Push(ud)
	return 1
}

func checkShellCmd(L *lua.LState) (*ShellCmd, *lua.LUserData) {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*ShellCmd); ok {
		return v, ud
	}
	L.ArgError(1, "shellcmd expected")
	return nil, nil
}

func getSetPath(L *lua.LState) int {
	c, ud := checkShellCmd(L)
	if c == nil {
		return 0
	}
	if L.GetTop() == 2 {
		c.Gocmd.Path = L.CheckString(2)
		L.Push(ud)
		return 1
	}
	L.Push(lua.LString(c.Gocmd.Path))
	return 1
}

func getSetDir(L *lua.LState) int {
	c, ud := checkShellCmd(L)
	if c == nil {
		return 0
	}
	if L.GetTop() == 2 {
		c.Gocmd.Dir = L.CheckString(2)
		L.Push(ud)
		return 1
	}
	L.Push(lua.LString(c.Gocmd.Dir))
	return 1
}

func getSetArgs(L *lua.LState) int {
	c, ud := checkShellCmd(L)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		// TODO: Take table->array, not vararg
		args := make([]string, L.GetTop()-1)
		for i := 2; i <= L.GetTop(); i++ {
			args[i-2] = L.CheckString(i)
		}
		c.Gocmd.Args = args
		L.Push(ud)
		return 1
	}
	ret := L.NewTable()
	for _, arg := range c.Gocmd.Args {
		ret.Append(lua.LString(arg))
	}
	L.Push(ret)
	return 1
}

func getSetEnv(L *lua.LState) int {
	c, ud := checkShellCmd(L)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		// TODO: Take table->dict, not vararg
		env := make([]string, L.GetTop()-1)
		for i := 2; i <= L.GetTop(); i++ {
			env[i-2] = L.CheckString(i)
		}
		c.Gocmd.Env = env
		L.Push(ud)
		return 1
	}
	// TODO: Return table->dict, not table->array
	ret := L.NewTable()
	for _, arg := range c.Gocmd.Env {
		ret.Append(lua.LString(arg))
	}
	L.Push(ret)
	return 1
}
