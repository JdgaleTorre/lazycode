package agent

import "context"

type OpenCodeAgent struct {
	BaseAgent
}

func NewOpenCodeAgent(command string) *OpenCodeAgent {
	if command == "" {
		command = "opencode"
	}
	return &OpenCodeAgent{BaseAgent{name: "opencode", command: command}}
}

func (a *OpenCodeAgent) StartSession(ctx context.Context, opts SessionOpts) (Session, error) {
	return a.StartPTYSession(ctx, opts, "/exit")
}
