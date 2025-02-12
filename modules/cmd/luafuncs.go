package cmd

import (
	"fmt"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func getSetErrorPropagation(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		val := L.CheckBool(2)
		c.lock.Lock()
		c.ErrorPropagation = val
		c.lock.Unlock()

		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	val := lua.LBool(c.ErrorPropagation)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func getSetDir(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		val := L.CheckString(2)
		c.lock.Lock()
		c.gocmd.Dir = val
		c.lock.Unlock()

		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	val := lua.LString(c.gocmd.Dir)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func getSetCmd(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		argsL := L.CheckTable(2)
		if argsL == nil {
			return 0
		}
		argsLLen := argsL.MaxN()
		if argsLLen < 1 {
			L.RaiseError("cmd must have at least one argument (the process binary)")
			return 0
		}
		args := make([]string, 0, argsLLen)
		for i := 1; i <= argsLLen; i++ {
			args = append(args, lua.LVAsString(argsL.RawGetInt(i)))
		}

		c.lock.Lock()
		c.gocmd.Path = args[0]
		c.gocmd.Args = args
		c.lock.Unlock()

		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	ret := L.NewTable()
	for _, arg := range c.gocmd.Args {
		ret.Append(lua.LString(arg))
	}
	c.lock.RUnlock()
	L.Push(ret)
	return 1
}

func getSetEnv(L *lua.LState) int {
	c, ud := checkCmd(L, 1)
	if c == nil {
		return 0
	}
	if L.GetTop() >= 2 {
		envL := L.CheckTable(2)
		if envL == nil {
			return 0
		}

		env := make([]string, 0)
		envK, envV := envL.Next(lua.LNil)
		for envK != lua.LNil {
			env = append(env, fmt.Sprintf("%s=%s", lua.LVAsString(envK), lua.LVAsString(envV)))
			envK, envV = envL.Next(envK)
		}

		c.lock.Lock()
		c.gocmd.Env = env
		c.lock.Unlock()

		L.Push(ud)
		return 1
	}

	c.lock.RLock()
	ret := L.NewTable()
	for _, arg := range c.gocmd.Env {
		envSplit := strings.SplitN(arg, "=", 2)
		if len(envSplit) < 2 {
			continue
		}
		ret.RawSetString(envSplit[0], lua.LString(envSplit[1]))
	}
	c.lock.RUnlock()
	L.Push(ret)
	return 1
}
