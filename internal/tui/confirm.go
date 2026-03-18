package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type confirmModel struct {
	name string
}

func newConfirmModel(name string) confirmModel {
	return confirmModel{name: name}
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (confirmModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			// Return confirm
			return m, nil
		case "n", "N", "esc":
			// Return cancel
			return m, nil
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintf("Are you sure you want to remove %q?\n\n", m.name))
	s.WriteString("(y/n)")
	return s.String()
}
