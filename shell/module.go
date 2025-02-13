package shell

const LuaName = "fox.shell"

func (s *Shell) Dependencies() []string {
	return []string{}
}

func (s *Shell) Interrupt(all bool) bool {
	cancelCtx := s.cancelCtx
	if cancelCtx != nil {
		cancelCtx()
		return true
	}
	return false
}

func (s *Shell) Name() string {
	return LuaName
}

func (s *Shell) PrePrompt() {
	// no-op
}
