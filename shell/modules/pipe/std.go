package pipe

import (
	"os"

	lua "github.com/yuin/gopher-lua"
)

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

func newStderrPipe(L *lua.LState) int {
	return pushPipe(L, &stderrPipe)
}

func newStdoutPipe(L *lua.LState) int {
	return pushPipe(L, &stdoutPipe)
}

func newStdinPipe(L *lua.LState) int {
	return pushPipe(L, &stdinPipe)
}

func newNullPipe(L *lua.LState) int {
	return pushPipe(L, &nullPipe)
}
