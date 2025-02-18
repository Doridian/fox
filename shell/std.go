package shell

import (
	"os"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/modules/pipe"
)

var osPipeCreator = pipe.FixedPipeCreator{
	Name: LuaName,
}

var stderrPipe = pipe.NewPipe(&osPipeCreator, "stderr", nil, os.Stderr, nil)
var stdoutPipe = pipe.NewPipe(&osPipeCreator, "stdout", nil, os.Stdout, nil)
var stdinPipe = pipe.NewPipe(&osPipeCreator, "stdin", os.Stdin, nil, nil)

func (s *Shell) Stdout() *pipe.Pipe {
	return stdoutPipe
}

func (s *Shell) Stderr() *pipe.Pipe {
	return stderrPipe
}

func (s *Shell) Stdin() *pipe.Pipe {
	return stdinPipe
}

func StdoutFor(loader *loader.LuaModule) *pipe.Pipe {
	shellMod := loader.GetModule(LuaName).(*Shell)
	return shellMod.Stdout()
}

func StderrFor(loader *loader.LuaModule) *pipe.Pipe {
	shellMod := loader.GetModule(LuaName).(*Shell)
	return shellMod.Stderr()
}

func StdinFor(loader *loader.LuaModule) *pipe.Pipe {
	shellMod := loader.GetModule(LuaName).(*Shell)
	return shellMod.Stdin()
}
