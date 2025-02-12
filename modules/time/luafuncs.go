package time

import (
	"time"

	foxtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

func doNow(L *lua.LState) int {
	return foxtime.Push(L, time.Now())
}
