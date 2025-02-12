package pipe

import (
	lua "github.com/yuin/gopher-lua"
)

func pipeClose(L *lua.LState) int {
	_, p, ud := Check(L, 1, false)
	if p == nil {
		return 0
	}

	p.Close()

	L.Push(ud)
	return 1
}

func pipeCanWrite(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	L.Push(lua.LBool(p.CanWrite()))
	return 1
}

func pipeCanRead(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	L.Push(lua.LBool(p.CanRead()))
	return 1
}

func pipeWrite(L *lua.LState) int {
	_, p, ud := Check(L, 1, false)
	if p == nil {
		return 0
	}
	data := L.CheckString(2)

	if p.wc == nil {
		if p.isNull {
			L.Push(ud)
			return 1
		}

		L.ArgError(1, "pipe must be a writer")
		return 0
	}

	_, err := p.wc.Write([]byte(data))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(ud)
	return 1
}

func pipeRead(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}
	len := int(L.CheckNumber(2))
	if len < 1 {
		L.ArgError(2, "len must be greater than 0")
		return 0
	}

	if p.rc == nil {
		if p.isNull {
			L.Push(lua.LString(""))
			return 1
		}

		L.ArgError(1, "pipe must be a reader")
		return 0
	}

	data := make([]byte, len)
	n, err := p.rc.Read(data)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(data[:n]))
	return 1
}

func pipeToString(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}

	L.Push(lua.LString(p.ToString()))
	return 1
}
