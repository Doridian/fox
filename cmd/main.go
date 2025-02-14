package main

import (
	"flag"

	"github.com/Doridian/fox/shell"

	_ "github.com/Doridian/fox/modules/loader/builtins"
)

func main() {
	flag.Parse()

	s := shell.New()

	runScript := flag.Arg(0)
	if runScript != "" {
		s.RunScript(runScript)
		return
	}

	running := true
	for running {
		running = s.RunPrompt()
	}
}
