package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type telepresenceMenuModel struct {
	list     list.Model
	item     item
	selected string
}

func newTelepresenceMenuModel(i item) telepresenceMenuModel {
	items := []list.Item{
		addMenuItem{title: "Connect", description: "Start Telepresence daemon and connect to cluster"},
		addMenuItem{title: "Disconnect", description: "Stop Telepresence daemon"},
		addMenuItem{title: "Restart", description: "Restart Telepresence connection"},
		addMenuItem{title: "Remove", description: "Remove Telepresence from managed tools"},
	}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Telepresence Actions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return telepresenceMenuModel{list: l, item: i}
}

func (m telepresenceMenuModel) Init() tea.Cmd {
	return nil
}

func (m telepresenceMenuModel) Update(msg tea.Msg) (telepresenceMenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(addMenuItem); ok {
				m.selected = i.title
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m telepresenceMenuModel) View() string {
	return m.list.View()
}
