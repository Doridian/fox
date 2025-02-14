package main

import (
	"flag"

	"github.com/Doridian/fox/shell"

	_ "github.com/Doridian/fox/modules/loader/builtins"
)

var continuePtr = flag.Bool("c", false, "Continue after running script")

func handleError() {

}

func main() {
	flag.Parse()

	s := shell.New()

	runScript := flag.Arg(0)
	if runScript != "" {
		s.RunFile(runScript)
		if !*continuePtr {
			return
		}
	}

	s.RunPrompt()
}
