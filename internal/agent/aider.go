package agent

import "context"

type AiderAgent struct {
	BaseAgent
}

func NewAiderAgent(command string) *AiderAgent {
	if command == "" {
		command = "aider"
	}
	return &AiderAgent{BaseAgent{name: "aider", command: command}}
}

func (a *AiderAgent) StartSession(ctx context.Context, opts SessionOpts) (Session, error) {
	return nil, nil
}
