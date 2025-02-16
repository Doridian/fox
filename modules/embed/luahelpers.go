package embed

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func (m *LuaModule) pushFileFromLua(L *lua.LState) int {
	name := L.CheckString(1)
	if name == "" {
		L.ArgError(1, "empty file name")
		return 0
	}

	fh, err := root.Open(name)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	defer fh.Close()

	lf, err := L.Load(fh, fmt.Sprintf("%s%s", m.getPrefix(L), name))
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(lf)
	return 1
}

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

func (m *LuaModule) getPrefix(L *lua.LState) string {
	return lua.LVAsString(L.GetField(m.mod, "prefix"))
}
