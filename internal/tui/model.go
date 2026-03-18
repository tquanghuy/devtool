package tui

import (
	"fmt"
	"os/exec"
	"strings"

	"devtool/internal/devtools"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)
)

type item struct {
	profile devtools.DevtoolProfile
}

func (i item) Title() string       { return i.profile.Name }
func (i item) Description() string { return fmt.Sprintf("%s %s", i.profile.Executable, i.profile.Args) }
func (i item) FilterValue() string { return i.profile.Name }

type state int

const (
	stateList state = iota
	stateAdd
	stateConfirm
	stateExec
)

type model struct {
	list         list.Model
	config       *devtools.DevtoolsConfig
	state        state
	addForm      addFormModel
	confirm      confirmModel
	err          error
	quitting     bool
}

func NewModel(cfg *devtools.DevtoolsConfig) model {
	var items []list.Item
	for _, p := range cfg.Devtools {
		items = append(items, item{profile: p})
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Developer Tools"
	l.Styles.Title = titleStyle
	
	// Add custom keybindings for 'a' and 'd' later, for now just standard list
	
	return model{
		list:   l,
		config: cfg,
		state:  stateList,
	}
}

type errMsg struct{ err error }

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.state == stateList {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "a":
				m.state = stateAdd
				m.addForm = newAddFormModel()
				return m, m.addForm.Init()
			case "d":
				if i, ok := m.list.SelectedItem().(item); ok {
					m.state = stateConfirm
					m.confirm = newConfirmModel(i.profile.Name)
					return m, nil
				}
			case "enter":
				if i, ok := m.list.SelectedItem().(item); ok {
					// Launch tool
					c := exec.Command(i.profile.Executable, strings.Fields(i.profile.Args)...)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return errMsg{err}
					})
				}
			}
		} else if m.state == stateAdd {
			if msg.String() == "esc" {
				m.state = stateList
				return m, nil
			}
			var cmd tea.Cmd
			m.addForm, cmd = m.addForm.Update(msg)
			if m.addForm.focused == fieldDone {
				// Form complete
				p := m.addForm.getProfile()
				if err := devtools.Add(m.config, p); err != nil {
					m.addForm.err = err
					m.addForm.focused = fieldArgs
					return m, nil
				}
				if err := devtools.Save(m.config); err != nil {
					m.addForm.err = err
					m.addForm.focused = fieldArgs
					return m, nil
				}
				// Refresh list
				m.list.InsertItem(len(m.config.Devtools)-1, item{profile: p})
				m.state = stateList
			}
			return m, cmd
		} else if m.state == stateConfirm {
			if msg.String() == "y" || msg.String() == "Y" {
				if i, ok := m.list.SelectedItem().(item); ok {
					if err := devtools.Remove(m.config, i.profile.Name); err == nil {
						_ = devtools.Save(m.config)
						m.list.RemoveItem(m.list.Index())
					}
				}
				m.state = stateList
				return m, nil
			} else if msg.String() == "n" || msg.String() == "N" || msg.String() == "esc" {
				m.state = stateList
				return m, nil
			}
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	if m.state == stateList {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}
	
	if len(m.config.Devtools) == 0 && m.state == stateList {
		return docStyle.Render("No devtools configured yet.\n\nPress 'a' to add your first tool, or 'q' to quit.")
	}

	switch m.state {
	case stateAdd:
		return docStyle.Render(m.addForm.View())
	case stateConfirm:
		return docStyle.Render(m.confirm.View())
	default:
		return docStyle.Render(m.list.View())
	}
}

// Start launches the interactive TUI.
func Start() error {
	cfg, err := devtools.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	p := tea.NewProgram(NewModel(cfg), tea.WithAltScreen())
	_, err = p.Run()
	return err
}
