package shell

import (
	"io"

	"github.com/Doridian/fox/modules/loader"
)

func (s *Shell) Stdout() io.Writer {
	return s.stdout
}

func (s *Shell) Stderr() io.Writer {
	return s.stderr
}

func (s *Shell) Stdin() io.Reader {
	return s.stdin
}

func StdoutFor(loader *loader.LuaModule) io.Writer {
	shellMod := loader.GetModule(LuaName).(*Shell)
	return shellMod.Stdout()
}

func StderrFor(loader *loader.LuaModule) io.Writer {
	shellMod := loader.GetModule(LuaName).(*Shell)
	return shellMod.Stderr()
}

func StdinFor(loader *loader.LuaModule) io.Reader {
	shellMod := loader.GetModule(LuaName).(*Shell)
	return shellMod.Stdin()
}
