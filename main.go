package main

import (
	"github.com/Doridian/fox/prompt"
	"github.com/Doridian/fox/shell"
)

func main() {
	p := prompt.NewPrompt()
	s := shell.NewShellManager()

	running := true
	for running {
		running = s.Run(p)
	}
}
