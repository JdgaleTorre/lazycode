package agent

import "context"

type ClaudeAgent struct {
	BaseAgent
}

func NewClaudeAgent(command string) *ClaudeAgent {
	if command == "" {
		command = "claude"
	}
	return &ClaudeAgent{BaseAgent{name: "claude", command: command}}
}

func (a *ClaudeAgent) StartSession(ctx context.Context, opts SessionOpts) (Session, error) {
	return a.StartPTYSession(ctx, opts, "/exit")
}
