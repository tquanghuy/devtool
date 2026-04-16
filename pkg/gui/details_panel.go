package gui

import (
	"fmt"
)

func (gui *Gui) updateDetails() {
	if !gui.showDetails {
		return
	}

	gui.details.Clear()

	focused := gui.app.GetFocus()
	if focused == gui.tools {
		gui.renderDetailsTool()
	} else if focused == gui.conns {
		gui.renderDetailsConn()
	} else if focused == gui.resources {
		gui.renderDetailsResource()
	} else {
		fmt.Fprint(gui.details, "[gray]Select a tool or connection to see details")
	}
}

func (gui *Gui) renderDetailsTool() {
	row, _ := gui.tools.GetSelection()
	if row <= 0 || row > len(gui.State.Tools) {
		fmt.Fprint(gui.details, "[gray]Select a tool to see details")
		return
	}

	tool := gui.State.Tools[row-1]

	fmt.Fprintf(gui.details, "[aqua]NAME: [white]%s\n", tool.Name)
	fmt.Fprintf(gui.details, "[aqua]KIND: [white]%s\n", tool.Kind)
	fmt.Fprint(gui.details, "\n[yellow]COMMANDS:\n")
	fmt.Fprintf(gui.details, " [green]Start: [white]%s\n", tool.StartCmd)
	fmt.Fprintf(gui.details, " [red]Stop:  [white]%s\n", tool.StopCmd)
	fmt.Fprintf(gui.details, " [blue]Check: [white]%s\n", tool.CheckCmd)

	if tool.DefaultPort != 0 {
		fmt.Fprintf(gui.details, "\n[aqua]DEFAULT PORT: [white]%d\n", tool.DefaultPort)
	}

	status := "STOPPED"
	if gui.OS.CheckToolStatus(tool.CheckCmd) {
		status = "RUNNING"
	}
	fmt.Fprintf(gui.details, "\n[aqua]STATUS: [white]%s\n", status)
}

func (gui *Gui) renderDetailsConn() {
	row, _ := gui.conns.GetSelection()
	if row <= 0 || row > len(gui.State.Connections) {
		fmt.Fprint(gui.details, "[gray]Select a connection to see details")
		return
	}

	conn := gui.State.Connections[row-1]

	fmt.Fprintf(gui.details, "[aqua]NAME:    [white]%s\n", conn.Name)
	fmt.Fprintf(gui.details, "[aqua]ADDRESS: [white]%s:%d\n", conn.Host, conn.Port)
	fmt.Fprintf(gui.details, "[aqua]HOST:    [white]%s\n", conn.Host)
	fmt.Fprintf(gui.details, "[aqua]PORT:    [white]%d\n", conn.Port)

	status := "DISCONNECTED"
	if gui.OS.DialTCP(conn.Host, conn.Port) {
		status = "CONNECTED"
	}
	fmt.Fprintf(gui.details, "\n[aqua]STATUS: [white]%s\n", status)
}

func (gui *Gui) renderDetailsResource() {
	row, _ := gui.resources.GetSelection()
	if row <= 0 || row > len(gui.State.Resources) {
		fmt.Fprint(gui.details, "[gray]Select a process to see details")
		return
	}

	proc := gui.State.Resources[row-1]

	fmt.Fprintf(gui.details, "[aqua]PROCESS: [white]%s\n", proc.Name)
	fmt.Fprintf(gui.details, "[aqua]PID:     [white]%s\n", proc.Type)
	fmt.Fprintf(gui.details, "[aqua]CPU:     [white]%s\n", proc.CPU)
	fmt.Fprintf(gui.details, "[aqua]MEMORY:  [white]%s\n", proc.MEM)
	fmt.Fprintf(gui.details, "[aqua]TYPE:    [white]%s\n", proc.Type)

	fmt.Fprint(gui.details, "\n[yellow]ACTIONS:\n")
	fmt.Fprint(gui.details, " [red]Enter: [white]Kill Process\n")
}
