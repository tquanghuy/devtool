package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type telepresenceInstallModel struct{}

func (m telepresenceInstallModel) Init() tea.Cmd {
	return nil
}

func (m telepresenceInstallModel) Update(msg tea.Msg) (telepresenceInstallModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "enter", "q":
			return m, nil // Handled by root
		}
	}
	return m, nil
}

func (m telepresenceInstallModel) View() string {
	return `Telepresence binary not found!

To manage Telepresence connections, please install it first:

macOS (Homebrew):
    brew install datawire/blackbird/telepresence

Linux/macOS (curl):
    sudo curl -fL https://app.getambassador.io/download/tel2/linux/amd64/latest/telepresence -o /usr/local/bin/telepresence
    sudo chmod a+x /usr/local/bin/telepresence

(Press enter or esc to go back)`
}
