package tui

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"devtool/internal/devtools"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

var (
	docStyle = lipgloss.NewStyle().Margin(1, 2)
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")).
			Bold(true).
			Padding(0, 1)
)

type item struct {
	profile devtools.DevtoolProfile
	status  TelepresenceStatus
}

func (i item) Title() string       { return i.profile.Name }
func (i item) Description() string {
	if i.profile.Name == "Telepresence" {
		statusStr := "Unknown"
		switch i.status {
		case StatusNotInstalled:
			statusStr = "Not Installed"
		case StatusDisconnected:
			statusStr = "Disconnected"
		case StatusConnected:
			statusStr = "Connected"
		}
		return fmt.Sprintf("Status: %s", statusStr)
	}
	return fmt.Sprintf("%s %s", i.profile.Executable, i.profile.Args)
}
func (i item) FilterValue() string { return i.profile.Name }

type state int

const (
	stateList state = iota
	stateAddMenu
	stateAdd
	stateConfirm
	stateExec
	stateTelepresenceMenu
	stateTelepresenceInstall
)

type model struct {
	list                list.Model
	config              *devtools.DevtoolsConfig
	state               state
	addMenu             addMenuModel
	addForm             addFormModel
	confirm             confirmModel
	telepresenceMenu    telepresenceMenuModel
	telepresenceInstall telepresenceInstallModel
	err                 error
	quitting            bool
}

func NewModel(cfg *devtools.DevtoolsConfig) model {
	var items []list.Item
	for _, p := range cfg.Devtools {
		items = append(items, item{profile: p, status: StatusUnknown})
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
	// Start polling status immediately if telepresence is in the list
	for _, it := range m.list.Items() {
		if i, ok := it.(item); ok && i.profile.Name == "Telepresence" {
			return CheckTelepresenceStatusCmd()
		}
	}
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case TelepresenceActionMsg:
		statusMsg := fmt.Sprintf("Telepresence %s completed", msg.Action)
		if msg.Err != nil {
			statusMsg = fmt.Sprintf("Telepresence %s failed: %v", msg.Action, msg.Err)
		}
		cmd := m.list.NewStatusMessage(statusMsg)
		return m, tea.Batch(CheckTelepresenceStatusCmd(), cmd)
	case TelepresenceStatusMsg:
		var items []list.Item
		for _, it := range m.list.Items() {
			if i, ok := it.(item); ok {
				if i.profile.Name == "Telepresence" {
					i.status = msg.Status
				}
				items = append(items, i)
			}
		}
		m.list.SetItems(items)
		
		// Poll again periodically
		return m, tea.Tick(5*time.Second, func(t time.Time) tea.Msg { return tickMsg(t) })
	case tickMsg:
		return m, CheckTelepresenceStatusCmd()
	case tea.KeyMsg:
		if m.state == stateList {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "a":
				m.state = stateAddMenu
				m.addMenu = newAddMenuModel()
				return m, m.addMenu.Init()
			case "d":
				if i, ok := m.list.SelectedItem().(item); ok {
					m.state = stateConfirm
					m.confirm = newConfirmModel(i.profile.Name)
					return m, nil
				}
			case "enter":
				if i, ok := m.list.SelectedItem().(item); ok {
					if i.profile.Name == "Telepresence" {
						if !CheckTelepresenceInstalled() {
							m.state = stateTelepresenceInstall
							m.telepresenceInstall = telepresenceInstallModel{}
							return m, nil
						}
						m.state = stateTelepresenceMenu
						m.telepresenceMenu = newTelepresenceMenuModel(i)
						return m, m.telepresenceMenu.Init()
					}
					// Launch tool
					c := exec.Command(i.profile.Executable, strings.Fields(i.profile.Args)...)
					return m, tea.ExecProcess(c, func(err error) tea.Msg {
						return errMsg{err}
					})
				}
			}
		} else if m.state == stateAddMenu {
			if msg.String() == "esc" {
				m.state = stateList
				return m, nil
			}
			var cmd tea.Cmd
			m.addMenu, cmd = m.addMenu.Update(msg)
			if m.addMenu.selected != nil {
				if m.addMenu.selected.isCustom {
					m.state = stateAdd
					m.addForm = newAddFormModel()
					return m, m.addForm.Init()
				} else {
					p := m.addMenu.selected.profile
					if err := devtools.Add(m.config, p); err == nil {
						_ = devtools.Save(m.config)
						m.list.InsertItem(len(m.config.Devtools)-1, item{profile: p})
					}
					m.state = stateList
					return m, CheckTelepresenceStatusCmd()
				}
			}
			return m, cmd
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
		} else if m.state == stateTelepresenceMenu {
			if msg.String() == "esc" {
				m.state = stateList
				return m, nil
			}
			var cmd tea.Cmd
			m.telepresenceMenu, cmd = m.telepresenceMenu.Update(msg)
			if m.telepresenceMenu.selected != "" {
				action := m.telepresenceMenu.selected
				m.state = stateList
				switch action {
				case "Connect":
					return m, tea.Batch(ConnectTelepresenceCmd(), m.list.NewStatusMessage("Connecting to Telepresence..."))
				case "Disconnect":
					return m, tea.Batch(DisconnectTelepresenceCmd(), m.list.NewStatusMessage("Disconnecting from Telepresence..."))
				case "Restart":
					return m, tea.Batch(tea.Sequence(DisconnectTelepresenceCmd(), ConnectTelepresenceCmd()), m.list.NewStatusMessage("Restarting Telepresence..."))
				case "Remove":
					if err := devtools.Remove(m.config, "Telepresence"); err == nil {
						_ = devtools.Save(m.config)
						for idx, it := range m.list.Items() {
							if itemObj, ok := it.(item); ok && itemObj.profile.Name == "Telepresence" {
								m.list.RemoveItem(idx)
								break
							}
						}
					}
					return m, nil
				}
			}
			return m, cmd
		} else if m.state == stateTelepresenceInstall {
			if msg.String() == "esc" || msg.String() == "enter" || msg.String() == "q" {
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
	case stateAddMenu:
		return docStyle.Render(m.addMenu.View())
	case stateAdd:
		return docStyle.Render(m.addForm.View())
	case stateConfirm:
		return docStyle.Render(m.confirm.View())
	case stateTelepresenceMenu:
		return docStyle.Render(m.telepresenceMenu.View())
	case stateTelepresenceInstall:
		return docStyle.Render(m.telepresenceInstall.View())
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
