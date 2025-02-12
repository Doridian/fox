package embedded

import lua "github.com/yuin/gopher-lua"

func readFileFromLua(L *lua.LState) []byte {
	name := L.CheckString(1)
	if name == "" {
		return nil
	}

	data, err := root.ReadFile(name)
	if err != nil {
		L.Error(lua.LString(err.Error()), 0)
		return nil
	}
	return data
}
