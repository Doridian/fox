package direntry

import (
	"github.com/Doridian/fox/shell/modules/fs/fileinfo"
	lua "github.com/yuin/gopher-lua"
)

func deName(L *lua.LState) int {
	de, _ := Check(L, 1)
	if de == nil {
		return 0
	}

	L.Push(lua.LString(de.Name()))
	return 1
}

func deIsDir(L *lua.LState) int {
	de, _ := Check(L, 1)
	if de == nil {
		return 0
	}

	L.Push(lua.LBool(de.IsDir()))
	return 1
}

func deType(L *lua.LState) int {
	de, _ := Check(L, 1)
	if de == nil {
		return 0
	}

	L.Push(lua.LNumber(de.Type()))
	return 1
}

func deInfo(L *lua.LState) int {
	de, _ := Check(L, 1)
	if de == nil {
		return 0
	}

	fi, err := de.Info()
	if err != nil {
		L.RaiseError("%v", err.Error())
		return 0
	}
	return fileinfo.Push(L, fi)
}
