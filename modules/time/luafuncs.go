package time

import (
	"time"

	subtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

func doNow(L *lua.LState) int {
	return subtime.Push(L, time.Now())
}
