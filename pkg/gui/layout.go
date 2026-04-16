package gui

import (
	"fmt"
	"github.com/rivo/tview"
)

func (gui *Gui) layout() {
	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(gui.resources, 8, 1, false)
	
	content := tview.NewFlex().SetDirection(tview.FlexColumn)
	content.AddItem(gui.tools, 0, 1, true)
	content.AddItem(gui.conns, 0, 1, false)
	
	flex.AddItem(content, 0, 1, true)
	flex.AddItem(gui.status, 1, 0, false)
	
	gui.pages.AddAndSwitchToPage("main", flex, true)

	gui.renderTools()
	gui.renderConnections()
	gui.renderResources()
	gui.renderStatus()
}

func (gui *Gui) renderStatus() {
	gui.status.Clear()
	fmt.Fprint(gui.status, " [cyan]tab:[white] switch panel • [red]q:[white] quit • [green]a:[white] add • [red]d:[white] delete • [yellow]enter:[white] select")
}
