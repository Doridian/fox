package main

import (
	"log"
	"os"

	"github.com/Doridian/fox/prompt"
	"github.com/Doridian/fox/shell"
)

func main() {
	p := prompt.NewPrompt()
	c := shell.NewShellManager()
	for {
		res, err := p.Prompt("fox> ")
		if err != nil {
			os.Stdout.Write([]byte("\n"))
			log.Printf("Prompt aborted: %v", err)
			break
		}
		if res != "" {
			doExit, exitCode := c.Run(p, res)
			if exitCode != 0 {
				log.Printf("Exit code: %v", exitCode)
			}
			if doExit {
				break
			}
		}
	}
}
