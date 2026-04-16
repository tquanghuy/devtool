package gui

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
	"devtool/pkg/commands"
)

func (gui *Gui) renderConnections() {
	gui.conns.Clear()
	gui.conns.SetFixed(1, 0)

	// Header row
	headerName := tview.NewTableCell(" NAME ").
		SetTextColor(gui.Theme.HeaderFg).
		SetBackgroundColor(gui.Theme.HeaderBg).
		SetAttributes(tcell.AttrBold).
		SetSelectable(false)
	headerAddr := tview.NewTableCell(" ADDRESS ").
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

	gui.conns.SetCell(0, 0, headerName)
	gui.conns.SetCell(0, 1, headerAddr)
	gui.conns.SetCell(0, 2, headerStatus)

	row := 1
	
	// Core Databases
	dbs := []struct {
		name string
		host string
		port int
		tool string
	}{
		{"Postgres", gui.Config.Postgres.Host, gui.Config.Postgres.Port, "postgres"},
		{"MySQL", gui.Config.MySQL.Host, gui.Config.MySQL.Port, "mysql"},
	}

	for _, db := range dbs {
		status := gui.checkConnection(db.tool, db.host, db.port)
		gui.addConnectionRow(row, db.name, fmt.Sprintf("%s:%d", db.host, db.port), status)
		row++
	}

	// Custom Connections
	if len(gui.Config.Connections) > 0 {
		for name, conn := range gui.Config.Connections {
			status := gui.checkConnection("", conn.Host, conn.Port)
			gui.addConnectionRow(row, name, fmt.Sprintf("%s:%d", conn.Host, conn.Port), status)
			row++
		}
	}
}

func (gui *Gui) addConnectionRow(row int, name, address, status string) {
	color := gui.Theme.Disconnected
	if status == "CONNECTED" {
		color = gui.Theme.Connected
	}

	nameCell := tview.NewTableCell(" " + name).SetTextColor(tview.Styles.PrimaryTextColor).SetExpansion(1)
	addrCell := tview.NewTableCell(address).SetTextColor(tview.Styles.SecondaryTextColor).SetExpansion(1)
	statusCell := tview.NewTableCell(status + " ").SetAlign(tview.AlignRight).SetTextColor(color)

	gui.conns.SetCell(row, 0, nameCell)
	gui.conns.SetCell(row, 1, addrCell)
	gui.conns.SetCell(row, 2, statusCell)
}

func (gui *Gui) checkConnection(toolName string, host string, port int) string {
	if toolName != "" {
		if t, ok := commands.LookupTool(toolName); ok {
			checkCmd := fmt.Sprintf(t.CheckCmd, port)
			if gui.OS.CheckToolStatus(checkCmd) {
				return "CONNECTED"
			}
		}
	}
	
	// Fallback to TCP dial
	if gui.OS.DialTCP(host, port) {
		return "CONNECTED"
	}
	
	return "DISCONNECTED"
}
