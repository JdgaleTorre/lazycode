package ui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dirEntry struct {
	name  string
	path  string
	isGit bool
}

func (m *ProjectModal) enterBrowseMode() {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/"
	}
	m.mode = modeBrowse
	m.browseDir = home
	m.cursorIdx = 0
	m.scrollOff = 0
	m.refreshEntries()
}

func (m *ProjectModal) refreshEntries() {
	tmp, err := readDirEntries(m.browseDir)
	if err != nil {
		m.entries = nil
		return
	}
	parent := filepath.Dir(m.browseDir)
	if parent != m.browseDir {
		tmp = append([]dirEntry{{name: "..", path: parent}}, tmp...)
	}
	m.entries = tmp
	if m.cursorIdx >= len(m.entries) {
		m.cursorIdx = 0
	}
}

func (m *ProjectModal) updateBrowse(msg tea.KeyMsg) (*ProjectModal, tea.Cmd) {
	switch {
	case msg.String() == "esc":
		m.mode = modeRecent
		m.cursorIdx = 0
		return m, nil
	case msg.String() == "enter" || msg.String() == " ":
		if m.cursorIdx < 0 || m.cursorIdx >= len(m.entries) {
			return m, nil
		}
		entry := m.entries[m.cursorIdx]
		if entry.name == ".." {
			m.browseDir = entry.path
			m.cursorIdx = 0
			m.scrollOff = 0
			m.refreshEntries()
			return m, nil
		}
		if entry.isGit {
			m.visible = false
			m.mode = modeRecent
			return m, func() tea.Msg {
				return ProjectConfirmMsg{Dir: entry.path}
			}
		}
		m.browseDir = entry.path
		m.cursorIdx = 0
		m.scrollOff = 0
		m.refreshEntries()
		return m, nil
	case msg.String() == "up" || msg.String() == "k":
		if m.cursorIdx > 0 {
			m.cursorIdx--
		}
	case msg.String() == "down" || msg.String() == "j":
		if m.cursorIdx < len(m.entries)-1 {
			m.cursorIdx++
		}
	case msg.String() == "backspace":
		parent := filepath.Dir(m.browseDir)
		if parent != m.browseDir {
			m.browseDir = parent
			m.cursorIdx = 0
			m.scrollOff = 0
			m.refreshEntries()
		}
	}
	return m, nil
}

func (m *ProjectModal) viewBrowse() string {
	title := TitleStyle.Render("Switch Project")

	dirLine := MutedStyle.Render(m.browseDir)

	maxVisible := m.height - 10
	if maxVisible < 3 {
		maxVisible = 3
	}

	if m.cursorIdx < m.scrollOff {
		m.scrollOff = m.cursorIdx
	}
	if m.cursorIdx >= m.scrollOff+maxVisible {
		m.scrollOff = m.cursorIdx - maxVisible + 1
	}

	var itemLines []string
	end := m.scrollOff + maxVisible
	if end > len(m.entries) {
		end = len(m.entries)
	}
	for _, entry := range m.entries[m.scrollOff:end] {
		idx := m.scrollOff + len(itemLines)
		prefix := "  "
		style := lipgloss.NewStyle().Foreground(ColorText)
		if idx == m.cursorIdx {
			prefix = "▸ "
			style = lipgloss.NewStyle().Foreground(ColorPrimary).Bold(true)
		}
		var label string
		if entry.isGit {
			label = style.Render(prefix + "● " + entry.name + "/")
		} else {
			label = style.Render(prefix + entry.name + "/")
		}
		itemLines = append(itemLines, label)
	}

	if len(m.entries) > end {
		itemLines = append(itemLines, MutedStyle.Render("  "+strings.Repeat("─", 36)))
	}

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		dirLine,
		MutedStyle.Render(strings.Repeat("─", 38)),
		strings.Join(itemLines, "\n"),
		"",
		MutedStyle.Render("↑↓ navigate  Enter select  Esc back"),
	)

	box := BorderFocused.
		Width(ModalWidth).
		Padding(1, 2).
		Render(content)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box)
}

func readDirEntries(dir string) ([]dirEntry, error) {
	ents, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []dirEntry
	for _, e := range ents {
		if !e.IsDir() {
			continue
		}
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		_, gitErr := os.Stat(filepath.Join(full, ".git"))
		isGit := gitErr == nil
		result = append(result, dirEntry{
			name:  e.Name(),
			path:  full,
			isGit: isGit,
		})
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i].isGit != result[j].isGit {
			return result[i].isGit
		}
		return result[i].name < result[j].name
	})

	return result, nil
}
