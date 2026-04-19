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
	} else if focused == gui.resources {
		gui.renderDetailsResource()
	} else {
		fmt.Fprint(gui.details, "[gray]Select a tool or resource to see details")
	}
}

func (gui *Gui) renderDetailsTool() {
	row, _ := gui.tools.GetSelection()
	if row <= 0 || row > len(gui.State.Tools) {
		fmt.Fprint(gui.details, "[gray]Select a tool to see details")
		return
	}

	tool := gui.State.Tools[row-1]

	fmt.Fprintf(gui.details, "[aqua]IDENTIFIER: [white]%s\n", tool.Instance.Identifier)
	fmt.Fprintf(gui.details, "[aqua]TOOL TYPE:  [white]%s\n", tool.Instance.ToolName)
	fmt.Fprintf(gui.details, "[aqua]KIND:       [white]%s\n", tool.Definition.Kind)

	fmt.Fprint(gui.details, "\n[yellow]COMMANDS:\n")
	fmt.Fprintf(gui.details, " [green]Start: [white]%s\n", gui.OS.FormatCommand(tool.Definition.StartCmd, tool.Instance.Port))
	fmt.Fprintf(gui.details, " [red]Stop:  [white]%s\n", gui.OS.FormatCommand(tool.Definition.StopCmd, tool.Instance.Port))
	fmt.Fprintf(gui.details, " [blue]Check: [white]%s\n", gui.OS.FormatCommand(tool.Definition.CheckCmd, tool.Instance.Port))

	if tool.Instance.Port != 0 {
		fmt.Fprintf(gui.details, "\n[aqua]ASSIGNED PORT: [white]%d\n", tool.Instance.Port)
	}

	status := "STOPPED"
	checkCmd := gui.OS.FormatCommand(tool.Definition.CheckCmd, tool.Instance.Port)
	if gui.OS.CheckToolStatus(checkCmd) {
		status = "RUNNING"
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
