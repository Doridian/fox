package embedded

import (
	"embed"

	lua "github.com/yuin/gopher-lua"
)

//go:embed root/*
var root embed.FS

func luaBareLoader(L *lua.LState) int {
	return luaLoader(L, "")
}

func luaPrefixLoader(L *lua.LState) int {
	return luaLoader(L, LuaName)
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
		L.Error(lua.LString(err.Error()), 0)
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
		L.Error(lua.LString(err.Error()), 0)
		return 0
	}
	return 0
}
