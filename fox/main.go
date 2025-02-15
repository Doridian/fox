package fox

import (
	"flag"
	"os"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
)

var continuePtr = flag.Bool("k", false, "Keep running after script")
var gomodsGlobal = flag.Bool("gomods-global", true, "Register go modules as globals")
var gomodsAutoLoad = flag.Bool("gomods-auto-load", true, "Automatically load go modules")
var commandPtr = flag.String("c", "", "Command to run")

func Main() error {
	flag.Parse()

	cfg := loader.DefaultConfig()
	cfg.Global = *gomodsGlobal
	cfg.AutoLoad = *gomodsAutoLoad
	loader.SetDefaultConfig(cfg)

	s := shell.New(flag.Args())

	command := *commandPtr
	if command != "" {
		return s.RunCommand(command)
	}

	runScript := flag.Arg(0)
	if runScript != "" {
		err := s.RunFile(runScript)
		if !*continuePtr {
			return err
		}
	}

	return s.RunPrompt()
}

func MainWithExit() {
	if Main() != nil {
		os.Exit(1)
	}
}
