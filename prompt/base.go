package prompt

import (
	"bufio"
	"errors"
	"os"
)

type PromptManager struct {
	inScanner *bufio.Scanner
}

func NewPrompt() *PromptManager {
	return &PromptManager{
		inScanner: bufio.NewScanner(os.Stdin),
	}
}

func (p *PromptManager) Prompt(disp string) (string, error) {
	os.Stdout.WriteString(disp)
	os.Stdout.Sync()
	if !p.inScanner.Scan() {
		err := p.inScanner.Err()
		if err == nil {
			err = errors.New("stdin closed")
		}
		return "", err
	}
	return p.inScanner.Text(), nil
}
