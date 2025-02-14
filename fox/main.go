package fox

import (
	"flag"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
)

var continuePtr = flag.Bool("c", false, "Continue after running script")
var gomodsGlobal = flag.Bool("gomods-global", true, "Register go modules as globals")
var gomodsAutoLoad = flag.Bool("gomods-auto-load", true, "Automatically load go modules")

func Main() {
	flag.Parse()

	cfg := loader.DefaultConfig()
	cfg.Global = *gomodsGlobal
	cfg.AutoLoad = *gomodsAutoLoad
	loader.SetDefaultConfig(cfg)

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
