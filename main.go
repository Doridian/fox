package main

import (
	"github.com/Doridian/fox/shell"
)

func main() {
	s := shell.NewShell()

	running := true
	for running {
		running = s.Run()
	}
}
