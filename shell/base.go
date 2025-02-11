package shell

import (
	"errors"
	"log"
	"strings"

	_ "embed"

	"github.com/Doridian/fox/prompt"
	"github.com/Doridian/fox/shell/modules/shellcmd"
	lua "github.com/yuin/gopher-lua"
)

//go:embed init.lua
var initCode string

type ShellManager struct {
	l *lua.LState
}

func NewShellManager() *ShellManager {
	s := &ShellManager{
		l: lua.NewState(),
	}
	s.init()
	return s
}

func (s *ShellManager) init() {
	mod := shellcmd.New()
	mod.Init(s.l)
	s.l.DoString(initCode)
}

// if .. then .. end
// while .. do .. end
// repeat .. until ..
// for .. do .. end
// .. = ..
// ( .. )

var ErrNeedMore = errors.New("Need more input")

func (s *ShellManager) CommandToLua(cmd string) (string, error) {
	lua := strings.Builder{}
	lua.WriteString(cmd)
	return lua.String(), nil
}

// Return true to exit shell
func (s *ShellManager) Run(p *prompt.PromptManager, cmd string) (bool, int) {
	var luaCode string
	var err error
	for {
		luaCode, err = s.CommandToLua(cmd)
		if err == nil {
			break
		}
		if err != ErrNeedMore {
			log.Printf("Error parsing command: %v", err)
			return false, 0
		}
		cmdAdd, err := p.Prompt("f.x> ")
		if err != nil {
			log.Printf("Prompt aborted: %v", err)
			return true, 0
		}
		cmd += cmdAdd
	}

	s.l.SetGlobal("_LAST_DO_EXIT", lua.LBool(false))
	s.l.SetGlobal("_LAST_EXIT_CODE", lua.LNumber(0))

	err = s.l.DoString(luaCode)
	if err != nil {
		log.Printf("Error running command: %v", err)
		return false, 1
	}
	retC := s.l.GetTop()
	if retC > 0 {
		s.l.Pop(retC)
	}

	doExit := lua.LVAsBool(s.l.GetGlobal("_LAST_DO_EXIT"))
	exitCode := lua.LVAsNumber(s.l.GetGlobal("_LAST_EXIT_CODE"))
	return doExit, int(exitCode)
}
