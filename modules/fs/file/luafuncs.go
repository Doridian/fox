package file

import (
	"fmt"

	"github.com/Doridian/fox/modules/fs/fileinfo"
	lua "github.com/yuin/gopher-lua"
)

func fileStat(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	fi, err := f.Stat()
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return fileinfo.Push(L, fi)
}

func fileToString(L *lua.LState) int {
	f, _ := Check(L, 1)
	if f == nil {
		return 0
	}

	stat, err := f.Stat()
	if err != nil {
		L.Push(lua.LString(fmt.Sprintf("%s{?%s}", LuaType, err)))
		return 1
	}

	L.Push(lua.LString(fmt.Sprintf("%s{%s}", LuaType, stat.Name())))
	return 1
}
