package gui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"devtool/pkg/commands"
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

	var toolsToShow []commands.ToolDefinition
	// Ensure core singletons are included
	coreTools := []string{"docker", "telepresence"}
	for _, name := range coreTools {
		if t, ok := commands.LookupTool(name); ok {
			toolsToShow = append(toolsToShow, t)
		}
	}

	// Add managed instances
	for _, inst := range gui.Managed.Instances {
		if t, ok := commands.LookupTool(inst.ToolName); ok {
			displayTool := t
			displayTool.Name = inst.Identifier
			toolsToShow = append(toolsToShow, displayTool)
		}
	}

	gui.State.Tools = toolsToShow

	for i, tool := range toolsToShow {
		statusText := "STOPPED"
		color := gui.Theme.Stopped
		if gui.OS.CheckToolStatus(tool.CheckCmd) {
			statusText = "RUNNING"
			color = gui.Theme.Running
		}
		
		nameCell := tview.NewTableCell(" " + tool.Name).
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
	toolNames := commands.GetSupportedToolNames()
	
	list := tview.NewList()
	for _, name := range toolNames {
		list.AddItem(name, "", 0, nil)
	}
	
	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		toolName := mainText
		err := gui.Managed.AddInstance(&config.ManagedInstance{
			ToolName:   toolName,
			Identifier: toolName,
		})
		if err == nil {
			config.SaveManagedConfig(gui.Managed)
			gui.renderTools()
		}
		gui.pages.RemovePage("addTool")
		gui.app.SetFocus(gui.tools)
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

func (gui *Gui) handleDeleteTool() {
	idx, _ := gui.tools.GetSelection()
	if idx < 0 || idx >= len(gui.State.Tools) {
		return
	}
	tool := gui.State.Tools[idx]
	
	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to delete %s?", tool.Name)).
		AddButtons([]string{"Delete", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Delete" {
				gui.Managed.RemoveInstance(tool.Name)
				config.SaveManagedConfig(gui.Managed)
				gui.renderTools()
			}
			gui.pages.RemovePage("deleteTool")
			gui.app.SetFocus(gui.tools)
		})
	
	gui.pages.AddPage("deleteTool", modal, true, true)
}

func (gui *Gui) handleToolEnter() {
	idx, _ := gui.tools.GetSelection()
	if idx < 0 || idx >= len(gui.State.Tools) {
		return
	}
	tool := gui.State.Tools[idx]

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Actions for %s", tool.Name)).
		AddButtons([]string{"Reconnect", "Disconnect", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			var err error
			switch buttonIndex {
			case 0: // Reconnect
				_, err = gui.OS.RunCommand(tool.StartCmd)
			case 1: // Disconnect
				_, err = gui.OS.RunCommand(tool.StopCmd)
			}
			if err == nil {
				gui.renderTools()
			}
			gui.pages.RemovePage("toolActions")
			gui.app.SetFocus(gui.tools)
		})

	gui.pages.AddPage("toolActions", modal, true, true)
}
