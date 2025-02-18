package cmd

import (
	"log"
	"sync"

	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/modules/vars"
	"github.com/Doridian/fox/shell"
	lua "github.com/yuin/gopher-lua"
)

const LuaName = "go:cmd"
const LuaTypeName = "Cmd"
const LuaType = LuaName + ":" + LuaTypeName

type LuaModule struct {
	allCmds    map[*Cmd]bool
	cmdRegLock sync.Mutex
	loader     *loader.LuaModule
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{
		allCmds: make(map[*Cmd]bool),
		loader:  loader,
	}
}

func (m *LuaModule) Loader(L *lua.LState) int {
	mod := L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"new":      m.newCmd,
		"run":      m.runCmd,
		"start":    m.startCmd,
		"lookPath": lookPath,

		"running": m.getRunning,
	})

	exitCodes := L.NewTable()
	exitCodes.RawSetString("LuaError", lua.LNumber(ExitCodeLuaError))
	exitCodes.RawSetString("ProcessCouldNotStart", lua.LNumber(ExitCodeProcessCouldNotStart))
	mod.RawSetString("ExitCodes", exitCodes)

	mt := L.NewTypeMetatable(LuaType)
	mt.RawSetString("__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"getDir": getDir,
		"dir":    setDir,
		"getCmd": getCmd,
		"cmd":    setCmd,
		"getEnv": getEnv,
		"env":    setEnv,

		"stdout":     setStdout,
		"getStdout":  getStdout,
		"stdoutPipe": acquireStdoutPipe,
		"stderr":     setStderr,
		"getStderr":  getStderr,
		"stderrPipe": acquireStderrPipe,
		"stdin":      setStdin,
		"getStdin":   getStdin,
		"stdinPipe":  acquireStdinPipe,

		"run":   doRun,
		"start": doStart,
		"wait":  doWait,
		"kill":  doKill,

		"getRaiseForBadExit": getRaiseForBadExit,
		"raiseForBadExit":    setRaiseForBadExit,
		"getAutoLookPath":    getAutoLookPath,
		"autoLookPath":       setAutoLookPath,
	}))
	L.SetFuncs(mt, map[string]lua.LGFunction{
		"__tostring": cmdToString,
		"__call":     doRun,
	})
	mod.RawSetString(LuaTypeName, mt)

	L.Push(mod)
	return 1
}

func (m *LuaModule) Dependencies() []string {
	return []string{shell.LuaName, vars.LuaName}
}

func (m *LuaModule) Name() string {
	return LuaName
}

func (m *LuaModule) Interrupt() bool {
	m.cmdRegLock.Lock()
	defer m.cmdRegLock.Unlock()

	triedKill := false
	for cmd := range m.allCmds {
		if !cmd.foreground {
			continue
		}

		if cmd.gocmd.ProcessState == nil && cmd.gocmd.Process != nil {
			_ = cmd.gocmd.Process.Kill()
			triedKill = true
		}
	}

	return triedKill
}

func (m *LuaModule) PrePrompt() {
	m.cmdRegLock.Lock()
	defer m.cmdRegLock.Unlock()

	toDelete := make([]*Cmd, 0)
	for cmd := range m.allCmds {
		exited := cmd.gocmd.Process == nil
		exitCode := 0
		if cmd.gocmd.ProcessState != nil {
			exited = true
			exitCode = cmd.gocmd.ProcessState.ExitCode()
		}

		if !exited {
			continue
		}
		toDelete = append(toDelete, cmd)

		if cmd.awaited {
			continue
		}

		if exitCode == 0 {
			log.Printf("job %s exited", cmd.ToString())
		} else {
			log.Printf("job %s exited with code %d", cmd.ToString(), exitCode)
		}
	}

	for _, cmd := range toDelete {
		delete(m.allCmds, cmd)
	}
}

func (m *LuaModule) addCmd(cmd *Cmd) {
	m.cmdRegLock.Lock()
	m.allCmds[cmd] = true
	m.cmdRegLock.Unlock()
}

func init() {
	loader.AddModuleDefault(newLuaModule)
}
