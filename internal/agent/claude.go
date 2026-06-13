package agent

import (
	"context"
	"os/exec"

	"github.com/josegale/lazycode/internal/terminal"
	"github.com/josegale/lazycode/internal/util"
)

type ClaudeAgent struct {
	command string
}

func NewClaudeAgent(command string) *ClaudeAgent {
	if command == "" {
		command = "claude"
	}
	return &ClaudeAgent{command: command}
}

func (a *ClaudeAgent) Name() string {
	return "claude"
}

func (a *ClaudeAgent) Available() (bool, error) {
	_, err := exec.LookPath(a.command)
	return err == nil, err
}

func (a *ClaudeAgent) StartSession(ctx context.Context, opts SessionOpts) (Session, error) {
	cmd := exec.CommandContext(ctx, a.command)
	cmd.Dir = opts.WorkDir

	ptyHandle, err := terminal.StartPTY(cmd)
	if err != nil {
		return nil, err
	}

	return NewPTYSession(ptyHandle, "claude", "claude", util.NewID(), "/exit"), nil
}

func (a *ClaudeAgent) ResumeSession(ctx context.Context, sessionID string) (Session, error) {
	// TODO: implement session resume
	return nil, nil
}
