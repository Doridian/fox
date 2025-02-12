package io

import (
	goio "io"

	lua "github.com/yuin/gopher-lua"
)

func ioClose(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	c, ok := f.(goio.Closer)
	if !ok {
		L.ArgError(1, "not closable")
		return 0
	}

	err := c.Close()
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return 0
}

func ioRead(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	r, ok := f.(goio.Reader)
	if !ok {
		L.ArgError(1, "not readable")
		return 0
	}

	len := int(L.CheckNumber(2))
	if len < 1 {
		L.ArgError(2, "len must be greater than 0")
		return 0
	}

	data := make([]byte, len)
	n, err := r.Read(data)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(data[:n]))
	return 1
}

func ioWrite(L *lua.LState) int {
	f, ud := Check(L, 1)
	if f == nil {
		return 0
	}

	w, ok := f.(goio.Writer)
	if !ok {
		L.ArgError(1, "not writable")
		return 0
	}

	data := L.CheckString(2)

	_, err := w.Write([]byte(data))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(ud)
	return 1
}

func ioSeek(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	s, ok := f.(goio.Seeker)
	if !ok {
		L.ArgError(1, "not seekable")
		return 0
	}

	offsetL := L.CheckNumber(2)
	whenceL := L.CheckNumber(3)

	n, err := s.Seek(int64(offsetL), int(whenceL))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LNumber(n))
	return 1
}
