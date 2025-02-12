package pipe

import lua "github.com/yuin/gopher-lua"

func pipeClose(L *lua.LState) int {
	ok, pipe, ud := CheckPipe(L, 1, false)
	if !ok {
		return 0
	}

	pipe.Close()

	L.Push(ud)
	return 1
}

func pipeCanWrite(L *lua.LState) int {
	ok, pipe, _ := CheckPipe(L, 1, false)
	if !ok {
		return 0
	}
	L.Push(lua.LBool(pipe.CanWrite()))
	return 1
}

func pipeCanRead(L *lua.LState) int {
	ok, pipe, _ := CheckPipe(L, 1, false)
	if !ok {
		return 0
	}
	L.Push(lua.LBool(pipe.CanRead()))
	return 1
}

func pipeWrite(L *lua.LState) int {
	ok, pipe, ud := CheckPipe(L, 1, false)
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

func pipeRead(L *lua.LState) int {
	ok, pipe, _ := CheckPipe(L, 1, false)
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
