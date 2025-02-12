package time

import (
	"time"

	modtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

func doNow(L *lua.LState) int {
	return modtime.Push(L, time.Now())
}
