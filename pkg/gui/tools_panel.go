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
	
	// Add core singletons from config if they are managed or should always be shown
	coreTools := []string{"docker", "telepresence", "postgres", "mysql"}
	for _, name := range coreTools {
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
		// Avoid double-adding core singletons
		isCore := false
		for _, ct := range coreTools {
			if inst.ToolName == ct {
				isCore = true
				break
			}
		}
		if isCore {
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
	identifier := toolName
	port := def.DefaultPort
	
	if def.Kind == config.PortBound {
		port = gui.OS.GetFreePort(def.DefaultPort)
		identifier = fmt.Sprintf("%s:%d", toolName, port)
	}

	form := tview.NewForm()
	form.AddInputField("Identifier", identifier, 30, nil, func(text string) {
		identifier = text
	})
	
	if def.Kind == config.PortBound {
		form.AddInputField("Port", fmt.Sprintf("%d", port), 10, nil, func(text string) {
			fmt.Sscanf(text, "%d", &port)
		})
	}

	form.AddButton("Add", func() {
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
		AddButtons([]string{"Reconnect", "Disconnect", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			var err error
			switch buttonIndex {
			case 0: // Reconnect
				formattedCmd := gui.OS.FormatCommand(tool.Definition.StartCmd, tool.Instance.Port)
				_, err = gui.OS.RunCommand(formattedCmd)
			case 1: // Disconnect
				formattedCmd := gui.OS.FormatCommand(tool.Definition.StopCmd, tool.Instance.Port)
				_, err = gui.OS.RunCommand(formattedCmd)
			}
			if err == nil {
				gui.renderTools()
			}
			gui.pages.RemovePage("toolActions")
			gui.app.SetFocus(gui.tools)
		})

	gui.pages.AddPage("toolActions", modal, true, true)
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
