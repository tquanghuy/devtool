package tui

import (
	"fmt"
	"strings"

	"devtool/internal/devtools"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type field int

const (
	fieldName field = iota
	fieldExec
	fieldArgs
	fieldDone
)

type addFormModel struct {
	inputs  []textinput.Model
	focused field
	err     error
}

func newAddFormModel() addFormModel {
	inputs := make([]textinput.Model, 3)
	
	inputs[fieldName] = textinput.New()
	inputs[fieldName].Placeholder = "Name (e.g. Postgres REPL)"
	inputs[fieldName].Focus()
	inputs[fieldName].CharLimit = 64
	inputs[fieldName].Width = 30

	inputs[fieldExec] = textinput.New()
	inputs[fieldExec].Placeholder = "Executable path (e.g. /usr/local/bin/psql)"
	inputs[fieldExec].CharLimit = 256
	inputs[fieldExec].Width = 50

	inputs[fieldArgs] = textinput.New()
	inputs[fieldArgs].Placeholder = "Default arguments (optional)"
	inputs[fieldArgs].CharLimit = 512
	inputs[fieldArgs].Width = 50

	return addFormModel{
		inputs:  inputs,
		focused: fieldName,
	}
}

func (m addFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m addFormModel) Update(msg tea.Msg) (addFormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, nil // Root model will handle return to menu

		case "enter":
			if m.focused == fieldArgs {
				// Validation check
				if m.inputs[fieldName].Value() == "" || m.inputs[fieldExec].Value() == "" {
					m.err = fmt.Errorf("name and executable are required")
					return m, nil
				}
				m.focused = fieldDone
				return m, nil
			}
			m.nextInput()
			return m, nil

		case "shift+tab", "up":
			m.prevInput()
			return m, nil

		case "tab", "down":
			m.nextInput()
			return m, nil
		}
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *addFormModel) nextInput() {
	m.inputs[m.focused].Blur()
	m.focused++
	if m.focused > fieldArgs {
		m.focused = fieldName
	}
	m.inputs[m.focused].Focus()
}

func (m *addFormModel) prevInput() {
	m.inputs[m.focused].Blur()
	m.focused--
	if m.focused < fieldName {
		m.focused = fieldArgs
	}
	m.inputs[m.focused].Focus()
}

func (m addFormModel) View() string {
	var s strings.Builder

	s.WriteString("Add New Developer Tool\n\n")

	for i := range m.inputs {
		s.WriteString(m.inputs[i].View())
		s.WriteString("\n")
	}

	if m.err != nil {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(fmt.Sprintf("\nError: %v", m.err)))
	}

	s.WriteString("\n(enter to continue, esc to cancel)\n")

	return s.String()
}

func (m addFormModel) getProfile() devtools.DevtoolProfile {
	return devtools.DevtoolProfile{
		Name:       m.inputs[fieldName].Value(),
		Executable: m.inputs[fieldExec].Value(),
		Args:       m.inputs[fieldArgs].Value(),
	}
}
