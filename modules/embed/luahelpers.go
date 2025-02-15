package embed

import (
	"fmt"
	"strings"

	"github.com/Doridian/fox/util"
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
		return util.RetNil(L)
	}
	mod := L.CheckTable(lua.UpvalueIndex(1))
	if mod == nil {
		return util.RetNil(L)
	}
	pathStr := lua.LVAsString(L.GetField(mod, "path"))
	if pathStr == "" {
		return util.RetNil(L)
	}

	fixedName := name
	if prefix != "" && !strings.HasPrefix(fixedName, prefix) {
		return util.RetNil(L)
	}
	fixedName = strings.ReplaceAll(fixedName, ".", "/")
	fixedName = strings.TrimPrefix(fixedName[len(prefix):], "/")

	errArr := []string{}

	paths := strings.Split(pathStr, ";")
	for _, path := range paths {
		fixedPath := strings.ReplaceAll(path, "?", fixedName)
		data, err := root.ReadFile(fixedPath)
		if err != nil {
			errArr = append(errArr, fmt.Sprintf("embed: module \"%s\" read error: %v", fixedPath, err))
			continue
		}
		lf, err := L.LoadString(string(data))
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
