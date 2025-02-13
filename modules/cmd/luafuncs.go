package cmd

import (
	"errors"
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
	err := c.setArgs(argsL)
	if err != nil {
		L.RaiseError("%v", err)
		return 0
	}

	L.Push(ud)
	return 1
}

func (c *Cmd) setArgs(argsL *lua.LTable) error {
	if argsL == nil || argsL == lua.LNil {
		c.lock.Lock()
		c.gocmd.Args = nil
		c.lock.Unlock()
		return nil
	}

	argsLLen := argsL.MaxN()
	if argsLLen < 1 {
		return errors.New("cmd must have at least one argument (the process binary)")
	}

	args := make([]string, 0, argsLLen)
	for i := 1; i <= argsLLen; i++ {
		args = append(args, lua.LVAsString(argsL.RawGetInt(i)))
	}

	c.lock.Lock()
	c.gocmd.Args = args
	c.lock.Unlock()
	return nil
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

	c.setEnv(envL)

	L.Push(ud)
	return 1
}

func (c *Cmd) setEnv(envL *lua.LTable) {
	if envL == nil || envL == lua.LNil {
		c.lock.Lock()
		c.gocmd.Env = nil
		c.lock.Unlock()
		return
	}

	env := make([]string, 0)
	envK, envV := envL.Next(lua.LNil)
	for envK != nil && envK != lua.LNil {
		env = append(env, fmt.Sprintf("%s=%s", lua.LVAsString(envK), lua.LVAsString(envV)))
		envK, envV = envL.Next(envK)
	}

	c.lock.Lock()
	c.gocmd.Env = env
	c.lock.Unlock()
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

func newCmdInt(L *lua.LState) (*Cmd, *lua.LUserData) {
	c := &Cmd{
		gocmd:        &exec.Cmd{},
		AutoLookPath: true,
	}

	// new([args, [dir, [env]]])
	c.setArgs(L.OptTable(1, nil))
	c.gocmd.Dir = L.OptString(2, "")
	c.setEnv(L.OptTable(3, nil))

	return c, ToUserdata(L, c)
}

func newCmd(L *lua.LState) int {
	_, ud := newCmdInt(L)
	L.Push(ud)
	return 1
}

func runCmd(L *lua.LState) int {
	c, ud := newCmdInt(L)
	if c == nil {
		return 0
	}

	return c.doRun(L, ud)
}

func startCmd(L *lua.LState) int {
	c, ud := newCmdInt(L)
	if c == nil {
		return 0
	}

	return c.doStart(L, ud)
}
