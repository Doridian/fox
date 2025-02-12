package embedded

import (
	"fmt"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func readFileFromLua(L *lua.LState) []byte {
	name := L.CheckString(1)
	if name == "" {
		L.ArgError(1, "empty file name")
		return nil
	}

	data, err := root.ReadFile(name)
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return nil
	}
	return data
}

func luaLoader(L *lua.LState, prefix string) int {
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

	if prefix != "" {
		fixedPrefix := strings.ReplaceAll(prefix, ".", "/")
		if !strings.HasSuffix(fixedPrefix, "/") {
			fixedPrefix += "/"
		}
		if !strings.HasPrefix(fixedName, fixedPrefix) {
			return 0
		}
		fixedName = fixedName[len(fixedPrefix):]
	}

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
