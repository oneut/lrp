package command

type NilCommand struct {
}

func (nc *NilCommand) Run(fn func(string)) {
}

func (nc *NilCommand) Start() {
}

func (nc *NilCommand) NeedsRestart() bool {
	return false
}

func (nc *NilCommand) Kill() bool {
	return false
}

func (nc *NilCommand) Stop() {
}
