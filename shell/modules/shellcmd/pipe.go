package shellcmd

import (
	"io"
	"os"

	lua "github.com/yuin/gopher-lua"
)

const luaShellPipeType = "shell/modules/shellcmd/Pipe"

var stderrPipe = Pipe{
	wc:           os.Stderr,
	forwardClose: false,
}

var stdoutPipe = Pipe{
	wc:           os.Stdout,
	forwardClose: false,
}

var stdinPipe = Pipe{
	rc:           os.Stdin,
	forwardClose: false,
}

var nullPipe = Pipe{
	isNull:       true,
	forwardClose: false,
}

type Pipe struct {
	cmd *ShellCmd

	forwardClose bool
	isNull       bool
	rc           io.ReadCloser
	wc           io.WriteCloser
}

func (p *Pipe) Close() {
	if !p.forwardClose {
		return
	}

	if p.rc != nil {
		p.rc.Close()
	}
	if p.wc != nil {
		p.wc.Close()
	}
}

func luaPipeClose(L *lua.LState) int {
	ok, pipe, ud := checkPipe(L, 1, false)
	if !ok {
		return 0
	}

	pipe.Close()

	L.Push(ud)
	return 1
}

func luaPipeWrite(L *lua.LState) int {
	ok, pipe, ud := checkPipe(L, 1, false)
	if !ok {
		return 0
	}
	data := L.CheckString(2)

	if pipe.wc == nil {
		if pipe.isNull {
			L.Push(ud)
			return 1
		}

		L.ArgError(1, "pipe must be a writer")
		return 0
	}

	_, err := pipe.wc.Write([]byte(data))
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	L.Push(ud)
	return 1
}

func luaPipeRead(L *lua.LState) int {
	ok, pipe, _ := checkPipe(L, 1, false)
	if !ok {
		return 0
	}
	len := int(L.CheckNumber(2))
	if len < 1 {
		L.ArgError(2, "len must be greater than 0")
		return 0
	}

	if pipe.rc == nil {
		if pipe.isNull {
			L.Push(lua.LString(""))
			return 1
		}

		L.ArgError(1, "pipe must be a reader")
		return 0
	}

	data := make([]byte, len)
	n, err := pipe.rc.Read(data)
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}

	L.Push(lua.LString(data[:n]))
	return 1
}

func newStderrPipe(L *lua.LState) int {
	return pushShellPipe(L, &stderrPipe)
}

func newStdoutPipe(L *lua.LState) int {
	return pushShellPipe(L, &stdoutPipe)
}

func newStdinPipe(L *lua.LState) int {
	return pushShellPipe(L, &stdinPipe)
}

func newNullPipe(L *lua.LState) int {
	return pushShellPipe(L, &nullPipe)
}
