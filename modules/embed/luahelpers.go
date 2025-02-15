package embed

import (
	"fmt"
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
		L.RaiseError("%v", err)
		return nil
	}
	return data
}

func luaLoaderInt(L *lua.LState, prefix string) int {
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
		L.Push(lua.LString("embed.path is not set"))
		return 1
	}

	fixedName := name
	if prefix != "" {
		if !strings.HasPrefix(fixedName, prefix) {
			L.Push(lua.LString(fmt.Sprintf("embed: module (%s) does not match prefix (%s)", fixedName, prefix)))
			return 1
		}
	}
	fixedName = strings.ReplaceAll(fixedName, ".", "/")
	fixedName = strings.TrimPrefix(fixedName[len(prefix):], "/")

	errStrBuilder := &strings.Builder{}

	paths := strings.Split(pathStr, ";")
	for _, path := range paths {
		fixedPath := strings.ReplaceAll(path, "?", fixedName)
		data, err := root.ReadFile(fixedPath)
		if err != nil {
			errStrBuilder.WriteString(fmt.Sprintf("embedded module %s read error: %v\n", fixedPath, err))
			continue
		}
		lf, err := L.LoadString(string(data))
		if err != nil {
			errStrBuilder.WriteString(fmt.Sprintf("embedded module %s load error: %v\n", fixedPath, err))
			continue
		}
		L.Push(lf)
		return 1
	}

	errStr := errStrBuilder.String()
	if len(errStr) > 0 {
		L.Push(lua.LString(errStr))
		return 1
	}

	return 0
}
