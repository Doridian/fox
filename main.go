package main

import (
	"log"

	"github.com/Doridian/fox/prompt"
	"github.com/Doridian/fox/shellcmd"
)

func main() {
	p := prompt.NewPrompt()
	c := shellcmd.NewShell()
	for {
		res, err := p.Prompt("fox> ")
		if err != nil {
			log.Printf("\nPrompt aborted: %v", err)
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
