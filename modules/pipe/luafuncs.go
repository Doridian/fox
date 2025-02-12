package pipe

import (
	"github.com/Doridian/fox/luautil"
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

	if p.wc == nil {
		if p.isNull {
			L.Push(ud)
			return 1
		}

		L.ArgError(1, "pipe must be a writer")
		return 0
	}

	luautil.WriterWrite(L, p.wc)
	L.Push(ud)
	return 1
}

func pipeRead(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
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

	return luautil.ReaderRead(L, p.rc)
}

func pipeToString(L *lua.LState) int {
	_, p, _ := Check(L, 1, false)
	if p == nil {
		return 0
	}

	L.Push(lua.LString(p.ToString()))
	return 1
}
