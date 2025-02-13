package time

import (
	"time"

	"github.com/Doridian/fox/modules/time/duration"
	subtime "github.com/Doridian/fox/modules/time/time"
	lua "github.com/yuin/gopher-lua"
)

func doNow(L *lua.LState) int {
	return subtime.Push(L, time.Now())
}

func timeSince(L *lua.LState) int {
	t, _ := subtime.Check(L, 1)
	return duration.Push(L, time.Since(t))
}

func timeUntil(L *lua.LState) int {
	t, _ := subtime.Check(L, 1)
	return duration.Push(L, time.Until(t))
}
