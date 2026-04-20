package gui

import (
	"fmt"
	"time"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"devtool/pkg/config"
)


func (gui *Gui) renderTools() {
	gui.tools.Clear()
	gui.tools.SetFixed(1, 0)

	// Header row
	headerName := tview.NewTableCell(" TOOL ").
		SetTextColor(gui.Theme.HeaderFg).
		SetBackgroundColor(gui.Theme.HeaderBg).
		SetAttributes(tcell.AttrBold).
		SetSelectable(false)
	headerStatus := tview.NewTableCell(" STATUS ").
		SetTextColor(gui.Theme.HeaderFg).
		SetBackgroundColor(gui.Theme.HeaderBg).
		SetAttributes(tcell.AttrBold).
		SetAlign(tview.AlignRight).
		SetSelectable(false)

	gui.tools.SetCell(0, 0, headerName)
	gui.tools.SetCell(0, 1, headerStatus)

	var toolsToShow []ToolInstance
	
	// Add core DB managers (always shown for now, though they don't have instances in 'Managed' usually)
	dbManagers := []string{"postgres", "mysql"}
	for _, name := range dbManagers {
		if t, ok := gui.Config.Tools[name]; ok {
			toolsToShow = append(toolsToShow, ToolInstance{
				Definition: t,
				Instance: config.ManagedInstance{
					ToolName:   name,
					Identifier: name,
					Port:       t.DefaultPort,
				},
			})
		}
	}

	// Add managed instances
	for _, inst := range gui.Config.Managed {
		// Avoid double-adding if it's already in toolsToShow (e.g. if we somehow managed a DB tool)
		exists := false
		for _, ts := range toolsToShow {
			if ts.Instance.Identifier == inst.Identifier {
				exists = true
				break
			}
		}
		if exists {
			continue
		}

		if t, ok := gui.Config.Tools[inst.ToolName]; ok {
			toolsToShow = append(toolsToShow, ToolInstance{
				Definition: t,
				Instance:   inst,
			})
		}
	}

	gui.State.Tools = toolsToShow

	for i, tool := range toolsToShow {
		statusText := "STOPPED"
		color := gui.Theme.Stopped
		
		// For DB tools, status is based on ANY connection being active? 
		// Or just show total connections count?
		// For now, keep it simple or show "TOOLS"
		if tool.Instance.ToolName == "postgres" || tool.Instance.ToolName == "mysql" {
			count := 0
			if tool.Instance.ToolName == "postgres" {
				count = len(gui.Config.PostgresConns)
			} else {
				count = len(gui.Config.MySQLConns)
			}
			statusText = fmt.Sprintf("%d CONNS", count)
			color = gui.Theme.HeaderBg
		} else if status, ok := gui.State.ToolStatuses[tool.Instance.Identifier]; ok {
			statusText = status
			if status == "RUNNING" {
				color = gui.Theme.Running
			}
		}
		
		nameCell := tview.NewTableCell(" " + tool.Instance.Identifier).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetExpansion(1)
		
		statusCell := tview.NewTableCell(statusText + " ").
			SetAlign(tview.AlignRight).
			SetTextColor(color)
		
		gui.tools.SetCell(i+1, 0, nameCell)
		gui.tools.SetCell(i+1, 1, statusCell)
	}
	
	gui.tools.SetSelectedFunc(func(row, column int) {
		if row == 0 {
			return
		}
		gui.handleToolEnter()
	})
}

func (gui *Gui) handleAddTool() {
	var toolNames []string
	for name := range gui.Config.Tools {
		if name == "postgres" || name == "mysql" {
			continue // These are core now
		}
		toolNames = append(toolNames, name)
	}
	
	list := tview.NewList()
	for _, name := range toolNames {
		list.AddItem(name, "", 0, nil)
	}
	
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		toolName := mainText
		toolDef := gui.Config.Tools[toolName]
		
		gui.pages.RemovePage("addTool")
		gui.showAddToolForm(toolName, toolDef)
	})
	
	list.SetDoneFunc(func() {
		gui.pages.RemovePage("addTool")
		gui.app.SetFocus(gui.tools)
	})
	
	list.SetBorder(true).SetTitle(" Add Tool (Esc to cancel) ")
	
	// Center the list
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 10, 1, true).
			AddItem(nil, 0, 1, false), 40, 1, true).
		AddItem(nil, 0, 1, false)

	gui.pages.AddPage("addTool", flex, true, true)
}

func (gui *Gui) showAddToolForm(toolName string, def config.ToolDefinition) {
	port := def.DefaultPort
	if port > 0 {
		port = gui.OS.GetFreePort(def.DefaultPort)
	}

	// Dynamic unique identifier generation
	identifier := toolName
	if port > 0 {
		identifier = fmt.Sprintf("%s:%d", toolName, port)
	}

	// Collision loop for identifier
	baseIdentifier := identifier
	counter := 1
	for {
		exists := false
		for _, inst := range gui.Config.Managed {
			if inst.Identifier == identifier {
				exists = true
				break
			}
		}
		if !exists {
			break
		}
		// If collision, add/increment suffix
		identifier = fmt.Sprintf("%s-%d", baseIdentifier, counter)
		counter++
	}

	form := tview.NewForm()
	form.AddInputField("Identifier", identifier, 40, nil, func(text string) {
		identifier = text
	})
	
	if def.DefaultPort > 0 {
		form.AddInputField("Port", fmt.Sprintf("%d", port), 10, nil, func(text string) {
			fmt.Sscanf(text, "%d", &port)
		})
	}

	form.AddButton("Add", func() {
		// Final uniqueness check before saving
		for _, inst := range gui.Config.Managed {
			if inst.Identifier == identifier {
				// TODO: Show a toast error instead of just returning
				return 
			}
		}

		gui.Config.Managed = append(gui.Config.Managed, config.ManagedInstance{
			ToolName:   toolName,
			Identifier: identifier,
			Port:       port,
			CreatedAt:  time.Now(),
		})
		_ = config.Save(gui.Config)
		gui.renderTools()
		gui.pages.RemovePage("addToolForm")
		gui.app.SetFocus(gui.tools)
	})

	form.AddButton("Cancel", func() {
		gui.pages.RemovePage("addToolForm")
		gui.app.SetFocus(gui.tools)
	})

	form.SetBorder(true).SetTitle(fmt.Sprintf(" Configure %s ", toolName)).SetTitleAlign(tview.AlignLeft)

	// Center the form
	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 12, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)

	gui.pages.AddPage("addToolForm", flex, true, true)
}


func (gui *Gui) handleDeleteTool() {
	idx, _ := gui.tools.GetSelection()
	if idx <= 0 || idx > len(gui.State.Tools) {
		return
	}
	tool := gui.State.Tools[idx-1]
	
	if tool.Instance.ToolName == "postgres" || tool.Instance.ToolName == "mysql" || tool.Instance.ToolName == "docker" || tool.Instance.ToolName == "telepresence" {
		// Cannot delete core tools
		return
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete %s?", tool.Instance.Identifier)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				for i, inst := range gui.Config.Managed {
					if inst.Identifier == tool.Instance.Identifier {
						gui.Config.Managed = append(gui.Config.Managed[:i], gui.Config.Managed[i+1:]...)
						_ = config.Save(gui.Config)
						gui.renderTools()
						break
					}
				}
			}
			gui.pages.RemovePage("deleteTool")
			gui.app.SetFocus(gui.tools)
		})
	
	gui.pages.AddPage("deleteTool", modal, true, true)
}

func (gui *Gui) handleToolEnter() {
	idx, _ := gui.tools.GetSelection()
	if idx <= 0 || idx > len(gui.State.Tools) {
		return
	}
	tool := gui.State.Tools[idx-1]

	if tool.Instance.ToolName == "postgres" || tool.Instance.ToolName == "mysql" {
		gui.showDatabaseConnections(tool.Instance.ToolName)
		return
	}

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Actions for %s", tool.Instance.Identifier)).
		AddButtons([]string{"Start/Restart", "Stop", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			gui.pages.RemovePage("toolActions")
			
			switch buttonIndex {
			case 0: // Start/Restart
				gui.startToolWithJIT(tool)
			case 1: // Stop
				formattedCmd := gui.OS.FormatCommand(tool.Definition.StopCmd, tool.Instance.Port)
				_, _ = gui.OS.RunCommand(formattedCmd)
				gui.renderTools()
			}
			gui.app.SetFocus(gui.tools)
		})

	gui.pages.AddPage("toolActions", modal, true, true)
}

func (gui *Gui) startToolWithJIT(tool ToolInstance) {
	if tool.Instance.Port > 0 && gui.OS.IsPortBusy(tool.Instance.Port) {
		// Check if it's already running (our own checkCmd)
		checkCmd := gui.OS.FormatCommand(tool.Definition.CheckCmd, tool.Instance.Port)
		if gui.OS.CheckToolStatus(checkCmd) {
			// Already running, just restart?
			formattedCmd := gui.OS.FormatCommand(tool.Definition.StartCmd, tool.Instance.Port)
			_, _ = gui.OS.RunCommand(formattedCmd)
			gui.renderTools()
			return
		}

		// It's busy but not by us (or check command failed)
		gui.showPortBusyModal(tool)
		return
	}

	// Normal start
	formattedCmd := gui.OS.FormatCommand(tool.Definition.StartCmd, tool.Instance.Port)
	_, _ = gui.OS.RunCommand(formattedCmd)
	gui.renderTools()
}

func (gui *Gui) showPortBusyModal(tool ToolInstance) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Port %d is already in use by another process.\nWould you like to allocate a new free port?", tool.Instance.Port)).
		AddButtons([]string{"Re-allocate Port", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Re-allocate Port" {
				newPort := gui.OS.GetFreePort(tool.Instance.Port + 1)
				if newPort > 0 {
					// Update Managed Instance in config
					for i, inst := range gui.Config.Managed {
						if inst.Identifier == tool.Instance.Identifier {
							gui.Config.Managed[i].Port = newPort
							// Also update identifier if it was tool:port
							if inst.Identifier == fmt.Sprintf("%s:%d", inst.ToolName, inst.Port) {
								gui.Config.Managed[i].Identifier = fmt.Sprintf("%s:%d", inst.ToolName, newPort)
							}
							break
						}
					}
					_ = config.Save(gui.Config)
					
					// Restart with new port
					tool.Instance.Port = newPort // Local update for immediate use
					gui.startToolWithJIT(tool)
				}
			}
			gui.pages.RemovePage("portBusy")
			gui.app.SetFocus(gui.tools)
		})

	gui.pages.AddPage("portBusy", modal, true, true)
}

func (gui *Gui) showDatabaseConnections(dbType string) {
	list := tview.NewList()
	
	conns := gui.Config.PostgresConns
	if dbType == "mysql" {
		conns = gui.Config.MySQLConns
	}

	for name, conn := range conns {
		status := "DISCONNECTED"
		if s, ok := gui.State.ConnStatuses[name]; ok {
			status = s
		}
		list.AddItem(name, fmt.Sprintf("%s:%d (%s)", conn.Host, conn.Port, status), 0, nil)
	}

	list.SetDoneFunc(func() {
		gui.pages.RemovePage("dbConns")
		gui.app.SetFocus(gui.tools)
	})

	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'a' {
			gui.showAddConnectionForm(dbType)
			return nil
		}
		if event.Rune() == 'd' {
			idx := list.GetCurrentItem()
			if idx >= 0 {
				name, _ := list.GetItemText(idx)
				gui.handleDeleteConnection(dbType, name)
			}
			return nil
		}
		return event
	})

	list.SetBorder(true).SetTitle(fmt.Sprintf(" %s Connections (a: add, d: delete, esc: back) ", dbType))

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(list, 15, 1, true).
			AddItem(nil, 0, 1, false), 60, 1, true).
		AddItem(nil, 0, 1, false)

	gui.pages.AddPage("dbConns", flex, true, true)
}

func (gui *Gui) showAddConnectionForm(dbType string) {
	form := tview.NewForm()
	name := ""
	host := "localhost"
	port := 5432
	if dbType == "mysql" {
		port = 3306
	}
	user := "postgres"
	if dbType == "mysql" {
		user = "root"
	}
	database := "postgres"
	if dbType == "mysql" {
		database = "mysql"
	}

	form.AddInputField("Connection Name", "", 20, nil, func(text string) { name = text })
	form.AddInputField("Host", host, 20, nil, func(text string) { host = text })
	form.AddInputField("Port", fmt.Sprintf("%d", port), 10, nil, func(text string) { fmt.Sscanf(text, "%d", &port) })
	form.AddInputField("User", user, 20, nil, func(text string) { user = text })
	form.AddInputField("Database", database, 20, nil, func(text string) { database = text })

	form.AddButton("Save", func() {
		if name == "" { return }
		newConn := config.DatabaseConfig{Host: host, Port: port, User: user, Database: database}
		if dbType == "postgres" {
			gui.Config.PostgresConns[name] = newConn
		} else {
			gui.Config.MySQLConns[name] = newConn
		}
		_ = config.Save(gui.Config)
		gui.pages.RemovePage("addConnForm")
		gui.showDatabaseConnections(dbType) // Refresh
	})

	form.AddButton("Cancel", func() {
		gui.pages.RemovePage("addConnForm")
	})

	form.SetBorder(true).SetTitle(fmt.Sprintf(" Add %s Connection ", dbType))

	flex := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(form, 15, 1, true).
			AddItem(nil, 0, 1, false), 50, 1, true).
		AddItem(nil, 0, 1, false)

	gui.pages.AddPage("addConnForm", flex, true, true)
}

func (gui *Gui) handleDeleteConnection(dbType, name string) {
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Delete connection '%s'?", name)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				if dbType == "postgres" {
					delete(gui.Config.PostgresConns, name)
				} else {
					delete(gui.Config.MySQLConns, name)
				}
				_ = config.Save(gui.Config)
			}
			gui.pages.RemovePage("deleteConn")
			gui.showDatabaseConnections(dbType) // Refresh
		})
	gui.pages.AddPage("deleteConn", modal, true, true)
}
