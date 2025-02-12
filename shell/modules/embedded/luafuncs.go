package embedded

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func (m *LuaModule) luaLoader(L *lua.LState) int {
	name := L.CheckString(1)
	if name == "" {
		return 0
	}
	mod := L.CheckTable(lua.UpvalueIndex(1))
	if mod == nil {
		return 0
	}

	log.Printf("Load attempt %v (%v) %d\n", name, mod, L.GetTop())

	return 0
}
