package util

import (
	"io"

	lua "github.com/yuin/gopher-lua"
)

func WriterWrite(L *lua.LState, wc io.Writer) {
	data := L.CheckString(2)

	_, err := wc.Write([]byte(data))
	if err != nil {
		L.RaiseError("%v", err)
		return
	}
}

func ReaderRead(L *lua.LState, rc io.Reader) int {
	len := int(L.CheckNumber(2))
	if len < 1 {
		L.ArgError(2, "len must be greater than 0")
		return 0
	}

	data := make([]byte, len)
	n, err := rc.Read(data)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(data[:n]))
	return 1
}
