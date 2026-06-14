package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskViewModel struct {
	vp      viewport.Model
	content strings.Builder
	width   int
	height  int
	label   string
}

func NewTaskViewModel(label string) TaskViewModel {
	vp := viewport.New(0, 0)
	return TaskViewModel{vp: vp, label: label}
}

func (m TaskViewModel) SetSize(w, h int) TaskViewModel {
	m.width = w
	m.height = h
	m.vp.Width = w
	m.vp.Height = h
	return m
}

func (m TaskViewModel) Append(data string) TaskViewModel {
	m.content.WriteString(data)
	m.vp.SetContent(m.content.String())
	m.vp.GotoBottom()
	return m
}

func (m TaskViewModel) Content() string {
	return m.content.String()
}

func (m TaskViewModel) Update(msg tea.Msg) (TaskViewModel, tea.Cmd) {
	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m TaskViewModel) View() string {
	if m.content.Len() == 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
			MutedStyle.Render("Running "+m.label+"..."))
	}
	return m.vp.View()
}
