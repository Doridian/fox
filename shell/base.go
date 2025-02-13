package shell

import (
	"errors"
	"log"
	"os"
	"strings"

	_ "embed"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/prompt"
	lua "github.com/yuin/gopher-lua"
)

//go:embed init.lua
var initCode string

const FailoverPrompt = "?fox?>"

type Shell struct {
	l   *lua.LState
	mod *lua.LTable
}

func NewShell() *Shell {
	s := &Shell{
		l: lua.NewState(lua.Options{
			SkipOpenLibs: true,
		}),
	}
	s.luaInit()
	return s
}

func luaExit(L *lua.LState) int {
	exitCodeL := lua.LVAsNumber(L.CheckNumber(1))
	os.Exit(int(exitCodeL))
	return 0
}

func (s *Shell) luaInit() {
	lua.OpenBase(s.l)
	lua.OpenPackage(s.l)
	lua.OpenTable(s.l)
	// lua.OpenIo(s.l)
	// lua.OpenOs(s.l)
	lua.OpenString(s.l)
	lua.OpenMath(s.l)
	lua.OpenDebug(s.l)
	lua.OpenChannel(s.l)
	lua.OpenCoroutine(s.l)

	mod := s.l.RegisterModule("shell", map[string]lua.LGFunction{
		"exit": luaExit,
	}).(*lua.LTable)
	s.l.SetGlobal("shell", mod)
	s.mod = mod

	mainMod := loader.NewLuaModule()
	mainMod.Load(s.l)

	err := s.l.DoString(initCode)
	if err != nil {
		log.Fatalf("Error initializing shell: %v", err)
	}
}

// if .. then .. end
// while .. do .. end
// repeat .. until ..
// for .. do .. end
// .. = ..
// ( .. )

var ErrNeedMore = errors.New("Need more input")

func (s *Shell) CommandToLua(cmd string) (string, error) {
	shellParser := s.mod.RawGetString("parser")
	if shellParser == nil || shellParser == lua.LNil {
		if cmd[len(cmd)-1] == '\\' {
			return "", ErrNeedMore
		}
		return strings.ReplaceAll(cmd, "\\\n", "\n"), nil
	}

	s.l.Push(shellParser)
	s.l.Push(lua.LString(cmd))
	err := s.l.PCall(1, 1, nil)
	if err != nil {
		return "", err
	}
	parseRet := s.l.Get(-1)
	s.l.Pop(1)
	if parseRet == lua.LTrue {
		return "", ErrNeedMore
	}
	return lua.LVAsString(parseRet), nil
}

func (s *Shell) RenderPrompt(lineNo int) (string, error) {
	renderPrompt := s.mod.RawGetString("renderPrompt")
	if renderPrompt == nil || renderPrompt == lua.LNil {
		if lineNo < 2 {
			return "fox> ", nil
		}
		return "fo+> ", nil
	}

	s.l.Push(renderPrompt)
	s.l.Push(lua.LNumber(lineNo))
	err := s.l.PCall(1, 1, nil)
	if err != nil {
		return FailoverPrompt, err
	}
	parseRet := s.l.Get(-1)
	s.l.Pop(1)
	return lua.LVAsString(parseRet), nil
}

func (s *Shell) Run(p *prompt.PromptManager) bool {
	pStr, err := s.RenderPrompt(1)
	if err != nil {
		log.Printf("Prompt rendering error: %v", err)
	}
	res, err := p.Prompt(pStr)
	if err != nil {
		os.Stdout.Write([]byte("\n"))
		log.Printf("Prompt aborted: %v", err)
		return false
	}
	if res != "" {
		exitCode := s.runOne(p, res)
		if exitCode != 0 {
			log.Printf("Exit code: %v", exitCode)
		}
	}
	return true
}

// Return true to exit shell
func (s *Shell) runOne(p *prompt.PromptManager, cmd string) int {
	var luaCode string
	var err error
	cmdB := strings.Builder{}
	cmdB.WriteString(cmd)
	lineNo := 2
	for {
		luaCode, err = s.CommandToLua(cmdB.String())
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Error parsing command: %v", err)
			return 0
		}

		pStr, err := s.RenderPrompt(lineNo)
		if err != nil {
			log.Printf("Prompt rendering error: %v", err)
		}

		cmdAdd, err := p.Prompt(pStr)
		if err != nil {
			log.Printf("Prompt aborted: %v", err)
			os.Exit(0)
			return 0
		}
		cmdB.WriteRune('\n')
		cmdB.WriteString(cmdAdd)
		lineNo++
	}

	s.l.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))
	err = s.l.DoString(luaCode)
	exitCode := int(lua.LVAsNumber(s.l.GetGlobal("_LAST_EXIT_CODE")))

	if err != nil {
		if exitCode == 0 {
			exitCode = 1
		}
		log.Printf("Error running command (code %d): %v", exitCode, err)
		return exitCode
	}

	retC := s.l.GetTop()
	if retC > 0 {
		s.l.Pop(retC)
	}

	return exitCode
}
