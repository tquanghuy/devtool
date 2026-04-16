package gui

import (
	"fmt"
	"strings"

	"devtool/pkg/commands"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (gui *Gui) renderResources() {
	gui.resources.Clear()
	gui.resources.SetFixed(1, 0)

	// Update Title with System Stats
	if gui.State.TotalCPU != "" {
		title := fmt.Sprintf(" Resource Monitoring (CPU: %s, Mem: %s) ", gui.State.TotalCPU, gui.State.TotalMem)
		gui.resources.SetTitle(title)
	}

	// Header row
	headers := []string{" PROCESS ", " PID ", " CPU% ", " MEM "}
	for i, header := range headers {
		cell := tview.NewTableCell(header).
			SetTextColor(gui.Theme.HeaderFg).
			SetBackgroundColor(gui.Theme.HeaderBg).
			SetAttributes(tcell.AttrBold).
			SetSelectable(false)
		
		if i >= 2 {
			cell.SetAlign(tview.AlignRight)
		}
		gui.resources.SetCell(0, i, cell)
	}

	row := 1

	_, _, width, _ := gui.resources.GetInnerRect()
	maxNameWidth := width - 30
	if maxNameWidth < 10 {
		maxNameWidth = 50 // Default for small or uninitialized width
	}

	// Use resources from state (populated by background polling)
	for _, stat := range gui.State.Resources {
		gui.addResourceRow(row, stat, maxNameWidth)
		row++
	}

	gui.resources.SetSelectedFunc(func(row, column int) {
		gui.handleResourceEnter()
	})
}

func (gui *Gui) addResourceRow(row int, stat commands.ResourceStat, maxNameWidth int) {
	displayName := stat.Name
	if len(displayName) > maxNameWidth {
		displayName = displayName[:maxNameWidth-3] + "..."
	}

	nameCell := tview.NewTableCell(" " + displayName).SetTextColor(tview.Styles.PrimaryTextColor).SetExpansion(1)
	typeCell := tview.NewTableCell(stat.Type).SetTextColor(tview.Styles.SecondaryTextColor)
	cpuCell := tview.NewTableCell(stat.CPU + " ").SetAlign(tview.AlignRight).SetTextColor(tcell.ColorOrange)
	memCell := tview.NewTableCell(stat.MEM + " ").SetAlign(tview.AlignRight).SetTextColor(tcell.ColorSkyblue)

	gui.resources.SetCell(row, 0, nameCell)
	gui.resources.SetCell(row, 1, typeCell)
	gui.resources.SetCell(row, 2, cpuCell)
	gui.resources.SetCell(row, 3, memCell)
}

func (gui *Gui) handleResourceEnter() {
	row, _ := gui.resources.GetSelection()
	if row <= 0 {
		return
	}

	pid := strings.TrimSpace(gui.resources.GetCell(row, 1).Text)
	name := strings.TrimSpace(gui.resources.GetCell(row, 0).Text)

	modal := tview.NewModal().
		SetText(fmt.Sprintf("Are you sure you want to kill process %s (PID: %s)?", name, pid)).
		AddButtons([]string{"Kill", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Kill" {
				gui.OS.KillPID(pid)
			}
			gui.pages.RemovePage("killProcess")
			gui.app.SetFocus(gui.resources)
		})

	gui.pages.AddPage("killProcess", modal, true, true)
}
