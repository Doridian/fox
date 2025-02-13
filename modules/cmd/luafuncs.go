package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	lua "github.com/yuin/gopher-lua"
)

func getErrorPropagation(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := lua.LBool(c.ErrorPropagation)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func setErrorPropagation(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	val := L.CheckBool(2)
	c.lock.Lock()
	c.ErrorPropagation = val
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getAutoLookPath(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := lua.LBool(c.AutoLookPath)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func setAutoLookPath(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	val := L.CheckBool(2)
	c.lock.Lock()
	c.AutoLookPath = val
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getDir(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}

	c.lock.RLock()
	val := lua.LString(c.gocmd.Dir)
	c.lock.RUnlock()
	L.Push(val)
	return 1
}

func setDir(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

	val := L.CheckString(2)
	c.lock.Lock()
	c.gocmd.Dir = val
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getCmd(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
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

func setCmd(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

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
	c.gocmd.Args = args
	c.lock.Unlock()

	L.Push(ud)
	return 1
}

func getEnv(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
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

func setEnv(L *lua.LState) int {
	c, ud := Check(L, 1)
	if c == nil {
		return 0
	}

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

func cmdToString(L *lua.LState) int {
	c, _ := Check(L, 1)
	if c == nil {
		return 0
	}
	L.Push(lua.LString(c.ToString()))
	return 1
}

func lookPath(L *lua.LState) int {
	cmd := L.CheckString(1)
	path, err := exec.LookPath(cmd)
	if err != nil {
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	L.Push(lua.LString(path))
	return 1
}

func newCmd(L *lua.LState) int {
	c := &Cmd{
		gocmd:        &exec.Cmd{},
		AutoLookPath: true,
	}
	return PushNew(L, c)
}
