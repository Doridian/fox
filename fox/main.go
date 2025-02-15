package fox

import (
	"errors"
	"flag"
	"os"
	"strconv"

	"github.com/Doridian/fox/modules/loader"
	"github.com/Doridian/fox/shell"
)

const (
	EvalUnset = iota
	EvalAsFile
	EvalAsString
	EvalAsCommand
)

var evalType = EvalUnset

func setEvalType(s string, newType int) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	if !v {
		return errors.New("flag must be true")
	}
	if evalType != EvalUnset {
		return errors.New("eval type already set (only at most one of -c, -e, -f can be set)")
	}
	evalType = newType
	return nil
}

var continuePtr = flag.Bool("k", false, "Keep running after script")
var gomodsGlobal = flag.Bool("gomods-global", true, "Register go modules as globals")
var gomodsAutoLoad = flag.Bool("gomods-auto-load", true, "Automatically load go modules")

func Main() error {
	flag.BoolFunc("c", "First arg is an internal command", func(s string) error {
		return setEvalType(s, EvalAsCommand)
	})
	flag.BoolFunc("e", "First arg is code to evaluate", func(s string) error {
		return setEvalType(s, EvalAsString)
	})
	flag.BoolFunc("f", "First arg is file to run (default)", func(s string) error {
		return setEvalType(s, EvalAsFile)
	})

	var err error

	flag.Parse()

	cfg := loader.DefaultConfig()
	cfg.Global = *gomodsGlobal
	cfg.AutoLoad = *gomodsAutoLoad
	loader.SetDefaultConfig(cfg)

	args := []string{}
	if flag.NArg() > 1 {
		args = flag.Args()[1:]
	}
	s := shell.New(args)

	if flag.NArg() > 0 {
		arg0 := flag.Arg(0)
		switch evalType {
		case EvalAsString:
			err = s.RunString(arg0)
		case EvalAsCommand:
			err = s.RunCommand(arg0)
		case EvalUnset, EvalAsFile:
			err = s.RunFile(arg0)
		}
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
