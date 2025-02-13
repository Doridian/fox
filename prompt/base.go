package prompt

import (
	"sync"

	"github.com/ergochat/readline"
)

type PromptManager struct {
	lock sync.Mutex
	rl   *readline.Instance
}

func NewPrompt() *PromptManager {
	rl, err := readline.New("?fox?> ")
	if err != nil {
		panic(err)
	}
	return &PromptManager{
		rl: rl,
	}
}

func (p *PromptManager) Prompt(disp string) (string, error) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.rl.SetPrompt(disp)
	return p.rl.ReadLine()
}
