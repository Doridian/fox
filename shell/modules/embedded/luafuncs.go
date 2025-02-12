package embedded

import (
	"embed"
	"fmt"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

//go:embed root/*
var root embed.FS

func (m *LuaModule) luaLoader(L *lua.LState) int {
	name := L.CheckString(1)
	if name == "" {
		return 0
	}
	mod := L.CheckTable(lua.UpvalueIndex(1))
	if mod == nil {
		return 0
	}
	pathStr := lua.LVAsString(L.GetField(mod, "path"))
	if pathStr == "" {
		return 0
	}

	fixedName := strings.ReplaceAll(name, ".", "/")

	paths := strings.Split(pathStr, ";")
	for _, path := range paths {
		fixedPath := strings.ReplaceAll(path, "?", fixedName)
		data, err := root.ReadFile(fixedPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			L.Push(lua.LString(fmt.Sprintf("embedded module %s read error: %v", fixedPath, err)))
			return 1
		}
		lf, err := L.LoadString(string(data))
		if err != nil {
			L.Push(lua.LString(fmt.Sprintf("embedded module %s load error: %v", fixedPath, err)))
			return 1
		}
		L.Push(lf)
		return 1
	}

	return 0
}
