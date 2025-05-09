package time

import (
	"fmt"
	gotime "time"

	"github.com/Doridian/fox/modules/duration"
	lua "github.com/yuin/gopher-lua"
)

func timeUnix(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Unix()))
	return 1
}

func timeUnixMilli(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.UnixMilli()))
	return 1
}

func timeUnixMicro(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.UnixMicro()))
	return 1
}

func timeUnixNano(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.UnixNano()))
	return 1
}

func timeSecond(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Second()))
	return 1
}

func timeNanosecond(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Nanosecond()))
	return 1
}

func timeHour(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Hour()))
	return 1
}

func timeMinute(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Minute()))
	return 1
}

func timeDay(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Day()))
	return 1
}

func timeMonth(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Month()))
	return 1
}

func timeWeekday(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Weekday()))
	return 1
}

func timeISOWeek(L *lua.LState) int {
	t, _ := Check(L, 1)
	year, week := t.ISOWeek()
	L.Push(lua.LNumber(year))
	L.Push(lua.LNumber(week))
	return 2
}

func timeDate(L *lua.LState) int {
	t, _ := Check(L, 1)
	y, m, d := t.Date()
	L.Push(lua.LNumber(y))
	L.Push(lua.LNumber(m))
	L.Push(lua.LNumber(d))
	return 2
}

func timeClock(L *lua.LState) int {
	t, _ := Check(L, 1)
	h, m, s := t.Clock()
	L.Push(lua.LNumber(h))
	L.Push(lua.LNumber(m))
	L.Push(lua.LNumber(s))
	return 2
}

func timeYearDay(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.YearDay()))
	return 1
}

func timeYear(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LNumber(t.Year()))
	return 1
}

func timeBefore(L *lua.LState) int {
	t, _ := Check(L, 1)
	t2, _ := Check(L, 2)
	L.Push(lua.LBool(t.Before(t2)))
	return 1
}

func timeNotAfter(L *lua.LState) int {
	t, _ := Check(L, 1)
	t2, _ := Check(L, 2)
	L.Push(lua.LBool(!t.After(t2)))
	return 1
}

func timeUTC(L *lua.LState) int {
	t, _ := Check(L, 1)
	return PushNew(L, t.UTC())
}

func timeLocal(L *lua.LState) int {
	t, _ := Check(L, 1)
	return PushNew(L, t.Local())
}

func timeAddDate(L *lua.LState) int {
	t, _ := Check(L, 1)
	years := L.CheckInt(2)
	months := L.CheckInt(3)
	days := L.CheckInt(4)
	return PushNew(L, t.AddDate(years, months, days))
}

func timeAddDuration(L *lua.LState) int {
	t, _ := Check(L, 1)
	d, _ := duration.Check(L, 2)
	return PushNew(L, t.Add(d))
}

func timeSubDuration(L *lua.LState) int {
	t, _ := Check(L, 1)
	d, _ := duration.Check(L, 2)
	return PushNew(L, t.Add(-d))
}

func timeDelta(L *lua.LState) int {
	t, _ := Check(L, 1)
	t2, _ := Check(L, 2)
	return duration.PushNew(L, t.Sub(t2))
}

func timeEq(L *lua.LState) int {
	d, _ := Check(L, 1)
	d2, _ := Check(L, 2)
	L.Push(lua.LBool(d.Equal(d2)))
	return 1
}

func timeSleepUntil(L *lua.LState) int {
	t, _ := Check(L, 1)
	gotime.Sleep(gotime.Until(t))
	return 0
}

func luaString(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LString(t.String()))
	return 1
}

func luaToString(L *lua.LState) int {
	t, _ := Check(L, 1)
	L.Push(lua.LString(fmt.Sprintf("%s{%s}", LuaType, t.String())))
	return 1
}

func luaFormat(L *lua.LState) int {
	t, _ := Check(L, 1)
	fmtStr := L.CheckString(2)
	L.Push(lua.LString(fmt.Sprintf("%s{%s}", LuaType, t.Format(fmtStr))))
	return 1
}

func timeNow(L *lua.LState) int {
	return PushNew(L, gotime.Now())
}

func timeSince(L *lua.LState) int {
	t, _ := Check(L, 1)
	return duration.PushNew(L, gotime.Since(t))
}

func timeUntil(L *lua.LState) int {
	t, _ := Check(L, 1)
	return duration.PushNew(L, gotime.Until(t))
}

func timeParse(L *lua.LState) int {
	tStr := L.CheckString(1)
	tLayout := L.OptString(2, gotime.RFC3339)
	t, err := gotime.Parse(tLayout, tStr)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}
	return PushNew(L, t)
}
