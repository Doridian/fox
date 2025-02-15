package fox

import (
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
)

var runFunc func(string) error

func setRunFunc(strVal string, newFunc func(string) error) error {
	boolVal, err := strconv.ParseBool(strVal)
	if err != nil {
		return err
	}
	if !boolVal {
		return errors.New("flag must be true")
	}
	if runFunc != nil {
		return errors.New("First arg type already set (only at most one of -c, -e, -f, -s can be set)")
	}
	runFunc = newFunc
	return nil
}

var continuePtr = flag.Bool("k", false, "Keep running after command/code/file")
var gomodsGlobal = flag.Bool("gomods-global", true, "Register go modules as globals")
var gomodsAutoLoad = flag.Bool("gomods-auto-load", true, "Automatically load go modules")

var forceContinue = false

func shellRunNoop(_ string) error {
	forceContinue = true
	return nil
}

func Main() error {
	s := shell.New()

	flag.BoolFunc("c", "First arg is an internal command", func(val string) error {
		return setRunFunc(val, s.RunCommand)
	})
	flag.BoolFunc("e", "First arg is code to evaluate", func(val string) error {
		return setRunFunc(val, s.RunString)
	})
	flag.BoolFunc("f", "First arg is file to run (default)", func(val string) error {
		return setRunFunc(val, s.RunFile)
	})
	flag.BoolFunc("s", "First arg is just passed to a shell", func(val string) error {
		return setRunFunc(val, shellRunNoop)
	})

	var err error

	flag.Parse()

	cfg := loader.DefaultConfig()
	cfg.Global = *gomodsGlobal
	cfg.AutoLoad = *gomodsAutoLoad
	loader.SetDefaultConfig(cfg)

	s.Init(flag.Args())

	if flag.NArg() > 0 {
		if runFunc == nil {
			runFunc = s.RunFile
		}
		err = runFunc(flag.Arg(0))
		if !forceContinue && !*continuePtr {
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
