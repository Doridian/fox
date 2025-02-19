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

var allCmds = make(map[*Cmd]bool)
var cmdRegLock sync.Mutex
var interactiveModeNeeded = false

type LuaModule struct {
	loader *loader.LuaModule
}

func newLuaModule(loader *loader.LuaModule) modules.LuaModule {
	return &LuaModule{
		loader: loader,
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
		"getDir":  getDir,
		"dir":     setDir,
		"getArgs": getArgs,
		"args":    setArgs,
		"getEnv":  getEnv,
		"env":     setEnv,

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
	cmdRegLock.Lock()
	defer cmdRegLock.Unlock()

	triedKill := false
	for cmd := range allCmds {
		if !cmd.foreground {
			continue
		}

		if cmd.gocmd.ProcessState == nil && cmd.gocmd.Process != nil {
			cmd.Stop()
			triedKill = true
		}
	}

	return triedKill
}

func (m *LuaModule) PrePrompt() {
	interactiveModeNeeded = true
	cleanupProcs()
}

func cleanupProc(cmd *Cmd) {
	cmdRegLock.Lock()
	defer cmdRegLock.Unlock()

	if cleanupProcNoLock(cmd) {
		delete(allCmds, cmd)
	}
}

func cleanupProcNoLock(cmd *Cmd) bool {
	exited := cmd.gocmd.Process == nil
	exitCode := 0
	if cmd.gocmd.ProcessState != nil {
		exited = true
		exitCode = cmd.gocmd.ProcessState.ExitCode()
	}

	if cmd.iCmd != nil {
		exitCode = cmd.iExit
		exited = cmd.iDone
	}

	if !exited {
		return false
	}

	if cmd.awaited {
		return true
	}

	if exitCode == 0 {
		log.Printf("job %s exited", cmd.ToString())
	} else {
		log.Printf("job %s exited with code %d", cmd.ToString(), exitCode)
	}
	return true
}

func cleanupProcs() {
	cmdRegLock.Lock()
	defer cmdRegLock.Unlock()

	toDelete := make([]*Cmd, 0)
	for cmd := range allCmds {
		if cleanupProcNoLock(cmd) {
			toDelete = append(toDelete, cmd)
		}
	}

	for _, cmd := range toDelete {
		delete(allCmds, cmd)
	}
}

func addCmd(cmd *Cmd) {
	cmdRegLock.Lock()
	allCmds[cmd] = true
	cmdRegLock.Unlock()
}

func delCmd(cmd *Cmd) {
	cmdRegLock.Lock()
	delete(allCmds, cmd)
	cmdRegLock.Unlock()
}

func init() {
	loader.AddModuleDefault(newLuaModule)
}
