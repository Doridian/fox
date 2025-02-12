package embed

import (
	goembed "embed"

	"github.com/Doridian/fox/modules/fs/direntry"
	lua "github.com/yuin/gopher-lua"
)

//go:embed root/*
var root goembed.FS

func luaLoader(L *lua.LState) int {
	mod := L.CheckTable(lua.UpvalueIndex(1))
	if mod == nil {
		return 0
	}
	prefixStr := lua.LVAsString(L.GetField(mod, "prefix"))
	return luaLoaderInt(L, prefixStr)
}

func luaReadFile(L *lua.LState) int {
	data := readFileFromLua(L)
	if data == nil {
		return 0
	}

	L.Push(lua.LString(string(data)))
	return 1
}

func luaLoadFile(L *lua.LState) int {
	data := readFileFromLua(L)
	if data == nil {
		return 0
	}

	lf, err := L.LoadString(string(data))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	L.Push(lf)
	return 1
}

func luaDoFile(L *lua.LState) int {
	data := readFileFromLua(L)
	if data == nil {
		return 0
	}

	err := L.DoString(string(data))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return 0
}

func luaReadDir(L *lua.LState) int {
	name := L.CheckString(1)
	if name == "" {
		L.ArgError(1, "empty dir name")
		return 0
	}

	dirents, err := root.ReadDir(name)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	return direntry.PushArray(L, dirents)
}
