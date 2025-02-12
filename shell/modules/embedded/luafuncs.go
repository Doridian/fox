package embedded

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func luaLoader(L *lua.LState) int {
	mod := L.CheckString(1)
	if mod == "" {
		return 0
	}

	log.Printf("Load attempt %v\n", mod)

	return 0
}
