package fox

import (
	"flag"

	"github.com/Doridian/fox/shell"
)

var continuePtr = flag.Bool("c", false, "Continue after running script")

func Main() {
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
