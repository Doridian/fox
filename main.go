package main

import (
	"github.com/Doridian/fox/shell"
)

func main() {
	s := shell.New()

	running := true
	for running {
		running = s.Run()
	}
}
