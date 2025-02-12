package shellcmd

import (
	"os/exec"

	"github.com/Doridian/fox/shell/modules/pipe"
	lua "github.com/yuin/gopher-lua"
)

const luaShellCmdType = "FOX/shellcmd/Cmd"

type Cmd struct {
	stdout *pipe.Pipe
	stderr *pipe.Pipe
	stdin  *pipe.Pipe

	gocmd            *exec.Cmd
	ErrorPropagation bool
}

func newShellCmd(L *lua.LState) int {
	return pushCmd(L, &Cmd{
		gocmd:            &exec.Cmd{},
		ErrorPropagation: false,
	})
}
