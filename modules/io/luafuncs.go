package io

import (
	"errors"
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
		if errors.Is(err, goio.EOF) {
			L.Push(lua.LNil)
			return 1
		}
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lua.LString(data[:n]))
	return 1
}

func ioReadToEnd(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	r, ok := f.(goio.Reader)
	if !ok {
		L.ArgError(1, "not readable")
		return 0
	}

	newData := make([]byte, 4096)
	data := []byte{}
	nTotal := 0
	for {
		n, err := r.Read(newData)
		nTotal += n
		if err != nil {
			if errors.Is(err, goio.EOF) {
				break
			}
			L.RaiseError("%v", err)
			return 0
		}
		if n > 0 {
			data = append(data, newData[:n]...)
		}
	}

	L.Push(lua.LString(data))
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

func ioCanWrite(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	_, ok := f.(goio.Writer)
	L.Push(lua.LBool(ok))
	return 1
}

func ioCanRead(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	_, ok := f.(goio.Reader)
	L.Push(lua.LBool(ok))
	return 1
}

func ioPrint(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	w, ok := f.(goio.Writer)
	if !ok {
		L.ArgError(1, "not writable")
		return 0
	}

	if L.GetTop() < 2 {
		return 0
	}

	var err error
	for i := 2; i <= L.GetTop(); i++ {
		lvRaw := L.Get(i)
		lvStrL := L.ToStringMeta(lvRaw)
		lvStr := lua.LVAsString(lvStrL)

		if i > 2 {
			_, err = w.Write([]byte("\t"))
			if err != nil {
				break
			}
		}
		_, err = w.Write([]byte(lvStr))
		if err != nil {
			break
		}
	}
	if err == nil {
		_, err = w.Write([]byte("\n"))
	}
	if err != nil {
		L.RaiseError("%v", err)
	}
	return 0
}
