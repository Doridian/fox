package time

import (
	gotime "time"

	"github.com/Doridian/fox/modules/duration"
	lua "github.com/yuin/gopher-lua"
)

func timeNow(L *lua.LState) int {
	return Push(L, gotime.Now())
}

func timeSince(L *lua.LState) int {
	t, _ := Check(L, 1)
	return duration.Push(L, gotime.Since(t))
}

func timeUntil(L *lua.LState) int {
	t, _ := Check(L, 1)
	return duration.Push(L, gotime.Until(t))
}

func timeParse(L *lua.LState) int {
	tStr := L.CheckString(1)
	tLayout := L.OptString(2, gotime.RFC3339)
	t, err := gotime.Parse(tLayout, tStr)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return Push(L, t)
}
