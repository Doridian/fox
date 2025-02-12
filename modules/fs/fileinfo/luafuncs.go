package fileinfo

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func fiName(L *lua.LState) int {
	fi, _ := Check(L, 1)
	if fi == nil {
		return 0
	}

	L.Push(lua.LString(fi.Name()))
	return 1
}

func fiSize(L *lua.LState) int {
	fi, _ := Check(L, 1)
	if fi == nil {
		return 0
	}

	L.Push(lua.LNumber(fi.Size()))
	return 1
}

func fiMode(L *lua.LState) int {
	fi, _ := Check(L, 1)
	if fi == nil {
		return 0
	}

	L.Push(lua.LNumber(fi.Mode()))
	return 1
}

func fiModTime(L *lua.LState) int {
	fi, _ := Check(L, 1)
	if fi == nil {
		return 0
	}

	L.Push(lua.LNumber(fi.ModTime().Unix()))
	return 1
}

func fiIsDir(L *lua.LState) int {
	fi, _ := Check(L, 1)
	if fi == nil {
		return 0
	}

	L.Push(lua.LBool(fi.IsDir()))
	return 1
}

func fiToString(L *lua.LState) int {
	fi, _ := Check(L, 1)
	if fi == nil {
		return 0
	}
	suffix := ""
	if fi.IsDir() {
		suffix = "/"
	}
	L.Push(lua.LString(fmt.Sprintf("%s{%s%s}", LuaType, fi.Name(), suffix)))
	return 1
}
