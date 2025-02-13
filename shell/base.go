package shell

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"

	_ "embed"

	"github.com/Doridian/fox/modules"
	"github.com/Doridian/fox/modules/cmd"
	"github.com/Doridian/fox/modules/loader"
	"github.com/ergochat/readline"
	lua "github.com/yuin/gopher-lua"
)

//go:embed init.lua
var initCode string

type Shell struct {
	l     *lua.LState
	lLock sync.Mutex

	mod   *lua.LTable
	print *lua.LFunction

	ctx       context.Context
	cancelCtx context.CancelFunc

	mainMod modules.LuaModule

	rlLock sync.Mutex
	rl     *readline.Instance

	signals chan os.Signal
}

// TODO: Handle SIGTERM

func NewShell() *Shell {
	rl, err := readline.New("?fox?> ")
	if err != nil {
		log.Panicf("Error initializing readline: %v", err)
	}

	s := &Shell{
		l: lua.NewState(lua.Options{
			SkipOpenLibs:        true,
			IncludeGoStackTrace: true,
		}),
		rl: rl,
	}
	s.luaInit()
	s.signalInit()
	return s
}

func luaExit(L *lua.LState) int {
	exitCodeL := lua.LVAsNumber(L.CheckNumber(1))
	os.Exit(int(exitCodeL))
	return 0
}

func (s *Shell) signalInit() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go func() {
		for range signals {
			interrupted := s.mainMod.Interrupt(false)
			if interrupted {
				continue
			}

			cancelCtx := s.cancelCtx
			if cancelCtx != nil {
				cancelCtx()
			}
		}
	}()
}

func (s *Shell) luaInit() {
	s.startLuaLock()
	defer s.endLuaLock()

	s.l.Pop(lua.OpenBase(s.l))
	s.l.Pop(lua.OpenPackage(s.l))
	s.l.Pop(lua.OpenTable(s.l))
	// s.l.Pop(lua.OpenIo(s.l))
	// s.l.Pop(lua.OpenOs(s.l))
	s.l.Pop(lua.OpenString(s.l))
	s.l.Pop(lua.OpenMath(s.l))
	s.l.Pop(lua.OpenDebug(s.l))
	s.l.Pop(lua.OpenChannel(s.l))
	s.l.Pop(lua.OpenCoroutine(s.l))

	mod := s.l.RegisterModule("shell", map[string]lua.LGFunction{
		"exit": luaExit,
	}).(*lua.LTable)
	s.l.SetGlobal("shell", mod)
	s.mod = mod

	s.print = s.l.GetGlobal("print").(*lua.LFunction)

	mainMod := loader.NewLuaModule()
	mainMod.Load(s.l)
	s.mainMod = mainMod

	err := s.l.DoString(initCode)
	if err != nil {
		log.Panicf("Error initializing shell: %v", err)
	}

	if s.l.GetTop() > 0 {
		log.Panicf("luaInit %d left stack frames!", s.l.GetTop())
	}
}

var ErrNeedMore = errors.New("Need more input")

func defaultShellParser(cmd string) (string, error) {
	if cmd[len(cmd)-1] == '\\' {
		return "", ErrNeedMore
	}
	return strings.ReplaceAll(cmd, "\\\n", "\n"), nil
}

func (s *Shell) shellParser(cmd string) (string, error) {
	shellParser := s.mod.RawGetString("parser")
	if shellParser == nil || shellParser == lua.LNil {
		return defaultShellParser(cmd)
	}

	s.startLuaLock()
	s.l.Push(shellParser)
	s.l.Push(lua.LString(cmd))
	err := s.l.PCall(1, 1, nil)
	s.endLuaLock()
	if err != nil {
		log.Printf("Error in Lua shell.parser: %v", err)
		return defaultShellParser(cmd)
	}
	parseRet := s.l.Get(-1)
	s.l.Pop(1)
	if parseRet == nil || parseRet == lua.LNil || parseRet == lua.LFalse {
		return defaultShellParser(cmd)
	} else if parseRet == lua.LTrue {
		return "", ErrNeedMore
	}
	return lua.LVAsString(parseRet), nil
}

func defaultRenderPrompt(lineNo int) string {
	if lineNo < 2 {
		return "fox> "
	}
	return "fo+> "
}

func (s *Shell) renderPrompt(lineNo int) string {
	renderPrompt := s.mod.RawGetString("renderPrompt")
	if renderPrompt == nil || renderPrompt == lua.LNil {
		return defaultRenderPrompt(lineNo)
	}

	s.startLuaLock()
	defer s.l.Push(renderPrompt)
	s.l.Push(lua.LNumber(lineNo))
	err := s.l.PCall(1, 1, nil)
	s.endLuaLock()
	if err != nil {
		log.Printf("Error in Lua shell.renderPrompt: %v", err)
		return defaultRenderPrompt(lineNo)
	}
	parseRet := s.l.Get(-1)
	s.l.Pop(1)

	if parseRet == nil || parseRet == lua.LNil || parseRet == lua.LFalse {
		return defaultRenderPrompt(lineNo)
	}
	return lua.LVAsString(parseRet)
}

func (s *Shell) readLine(disp string) (string, error) {
	s.rlLock.Lock()
	defer s.rlLock.Unlock()

	s.rl.SetPrompt(disp)
	return s.rl.ReadLine()
}

func (s *Shell) Run() bool {
	s.mainMod.PrePrompt()

	res, err := s.readLine(s.renderPrompt(1))
	if err != nil {
		if errors.Is(err, readline.ErrInterrupt) {
			return true
		}
		log.Printf("Prompt aborted: %v", err)
		return false
	}
	if res != "" {
		exitCode := s.runOne(res)
		if exitCode != 0 {
			log.Printf("Exit code: %v", exitCode)
		}
	}
	return true
}

func (s *Shell) startLuaLock() {
	s.lLock.Lock()
	s.ctx, s.cancelCtx = context.WithCancel(context.Background())
	s.l.SetContext(s.ctx)
}

func (s *Shell) endLuaLock() {
	cancelCtx := s.cancelCtx
	if cancelCtx != nil {
		cancelCtx()
	}
	s.ctx = nil
	s.cancelCtx = nil
	s.lLock.Unlock()
}

// Return true to exit shell
func (s *Shell) runOne(firstLine string) int {
	var luaCode string
	var err error

	cmdBuilder := strings.Builder{}
	cmdBuilder.WriteString(firstLine)

	lineNo := 1

	for {
		cmdBuilder.WriteRune('\n')

		luaCode, err = s.shellParser(cmdBuilder.String())
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Error parsing command: %v", err)
			return 0
		}

		lineNo++
		cmdAdd, err := s.readLine(s.renderPrompt(lineNo))
		if err != nil {
			if errors.Is(err, readline.ErrInterrupt) {
				return 0
			}
			log.Printf("Prompt aborted: %v", err)
			os.Exit(0)
			return 0
		}
		cmdBuilder.WriteString(cmdAdd)
	}

	if luaCode == "" {
		return 0
	}

	s.startLuaLock()
	defer s.endLuaLock()
	s.l.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))
	err = s.l.DoString(luaCode)
	exitCode := int(lua.LVAsNumber(s.l.GetGlobal("_LAST_EXIT_CODE")))

	if err != nil {
		if exitCode == 0 {
			exitCode = cmd.ExitCodeInternalShellError
		}
		log.Printf("Internal error running command: %v", err)
		return exitCode
	}

	retC := s.l.GetTop()
	if retC > 0 {
		s.l.Insert(s.print, 0)
		s.l.Call(retC, 0)
	}

	return exitCode
}
