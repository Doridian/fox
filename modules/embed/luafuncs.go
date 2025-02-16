package embed

import (
	goembed "embed"
	"fmt"
	"strings"

	"github.com/Doridian/fox/modules/fs/direntry"
	"github.com/Doridian/fox/modules/fs/file"
	"github.com/Doridian/fox/util"
	lua "github.com/yuin/gopher-lua"
)

//go:embed root/*
var root goembed.FS

func luaReadFile(L *lua.LState) int {
	data := readFileFromLua(L)
	if data == nil {
		return 0
	}

	L.Push(lua.LString(string(data)))
	return 1
}

func luaOpenFile(L *lua.LState) int {
	name := L.CheckString(1)
	if name == "" {
		L.ArgError(1, "empty file name")
		return 0
	}

	f, err := root.Open(name)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return file.PushNew(L, f)
}

func (m *LuaModule) luaLoadFile(L *lua.LState) int {
	return m.pushFileFromLua(L)
}

func (m *LuaModule) luaDoFile(L *lua.LState) int {
	m.pushFileFromLua(L)
	L.Call(0, 0)
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

func (m *LuaModule) luaLoader(L *lua.LState) int {
	prefix := m.getPrefix(L)

	name := L.CheckString(1)
	if name == "" {
		return util.RetNil(L)
	}
	pathStr := lua.LVAsString(L.GetField(m.mod, "path"))
	if pathStr == "" {
		return util.RetNil(L)
	}

	if prefix != "" && !strings.HasPrefix(name, prefix) {
		return util.RetNil(L)
	}
	fileName := name[len(prefix):]
	fileName = strings.ReplaceAll(fileName, ".", "/")
	fileName = strings.TrimPrefix(fileName, "/")

	errArr := []string{}

	paths := strings.Split(pathStr, ";")
	for _, path := range paths {
		fixedPath := strings.ReplaceAll(path, "?", fileName)
		fh, err := root.Open(fixedPath)
		if err != nil {
			errArr = append(errArr, fmt.Sprintf("embed: module \"%s\" read error: %v", fixedPath, err))
			continue
		}
		defer fh.Close()
		lf, err := L.Load(fh, name)
		if err != nil {
			errArr = append(errArr, fmt.Sprintf("embed: module \"%s\" load error: %v", fixedPath, err))
			continue
		}
		L.Push(lf)
		return 1
	}

	errStr := strings.Join(errArr, "\n\t")
	if len(errStr) > 0 {
		L.Push(lua.LString(errStr))
		return 1
	}

	return util.RetNil(L)
}
