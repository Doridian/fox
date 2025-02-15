package fox

import (
	"errors"
	"flag"
	"log"
	"os"
	"strconv"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
)

var runFunc func(string) error
var runFuncSet bool

func setRunFunc(strVal string, newFunc func(string) error) error {
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		return err
	}
	if !boolVal {
		return errors.New("flag must be true")
	}
	if runFuncSet {
		return errors.New("First arg type already set (only at most one of -c, -e, -f, -s can be set)")
	}
	runFuncSet = true
	runFunc = newFunc
	return nil
}

var continuePtr = flag.Bool("k", false, "Keep running after command/code/file")
var gomodsGlobal = flag.Bool("gomods-global", true, "Register go modules as globals")
var gomodsAutoload = flag.Bool("gomods-auto-load", true, "Automatically load go modules")

func shellRunNoop(_ string) error {
	return nil
}

func Main() error {
	var err error
	s := shell.New()

	forceShell := false
	flag.BoolFunc("c", "First arg is an internal command (default)", func(val string) error {
		return setRunFunc(val, s.RunCommand)
	})
	flag.BoolFunc("e", "First arg is code to evaluate", func(val string) error {
		return setRunFunc(val, s.RunString)
	})
	flag.BoolFunc("f", "First arg is file to run", func(val string) error {
		return setRunFunc(val, s.RunFile)
	})
	flag.BoolFunc("s", "First arg is just passed to a shell", func(val string) error {
		forceShell = true
		return setRunFunc(val, nil)
	})

	flag.Parse()

	cfg := loader.DefaultConfig()
	cfg.Global = *gomodsGlobal
	cfg.Autoload = *gomodsAutoload
	loader.SetDefaultConfig(cfg)

	if forceShell || flag.NArg() == 0 {
		if runFunc != nil {
			log.Fatalf("cannont run in non-shell mode without at least one argument")
		}
		s.Init(flag.Args())
		return s.RunPrompt()
	}

	if flag.NArg() > 1 {
		s.Init(flag.Args()[1:])
	} else {
		s.Init([]string{})
	}

	if flag.NArg() > 0 {
		if runFunc == nil {
			runFunc = s.RunCommand
		}
		err = runFunc(flag.Arg(0))
		if !forceShell && !*continuePtr {
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
