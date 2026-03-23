package tui

import (
	"devtool/internal/devtools"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type addMenuItem struct {
	title       string
	description string
	isCustom    bool
	profile     devtools.DevtoolProfile
}

func (i addMenuItem) Title() string       { return i.title }
func (i addMenuItem) Description() string { return i.description }
func (i addMenuItem) FilterValue() string { return i.title }

type addMenuModel struct {
	list     list.Model
	selected *addMenuItem
}

func newAddMenuModel() addMenuModel {
	items := []list.Item{
		addMenuItem{
			title:       "Custom Tool...",
			description: "Add a custom executable with arguments",
			isCustom:    true,
		},
		addMenuItem{
			title:       "Telepresence",
			description: "Manage Telepresence connection to cluster",
			isCustom:    false,
			profile: devtools.DevtoolProfile{
				Name:       "Telepresence",
				Executable: "telepresence",
				Args:       "", // Args aren't used natively as we will intercept this
			},
		},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select Tool to Add"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return addMenuModel{list: l}
}

func (m addMenuModel) Init() tea.Cmd {
	return nil
}

func (m addMenuModel) Update(msg tea.Msg) (addMenuModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if i, ok := m.list.SelectedItem().(addMenuItem); ok {
				m.selected = &i
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

func (m addMenuModel) View() string {
	return m.list.View()
}
