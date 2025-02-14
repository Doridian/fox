package main

import (
	"github.com/Doridian/fox/shell"

	_ "github.com/Doridian/fox/modules/loader/builtins"
)

func main() {
	s := shell.New()

	running := true
	for running {
		running = s.Run()
	}
}
