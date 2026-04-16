package gui

import (
	"github.com/gdamore/tcell/v2"
)

func (gui *Gui) keybindings() error {
	gui.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlC {
			gui.app.Stop()
			return nil
		}

		if event.Key() == tcell.KeyTab {
			gui.nextView()
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			name, _ := gui.pages.GetFrontPage()
			if name == "main" {
				gui.handleQuit()
				return nil
			}
		}

		if event.Rune() == 'a' {
			gui.handleAddTool()
			return nil
		}

		if event.Rune() == 'd' {
			gui.handleDeleteTool()
			return nil
		}

		if event.Rune() == 'i' {
			gui.showDetails = !gui.showDetails
			gui.layout()
			return nil
		}

		if event.Rune() == 'i' {
			gui.showDetails = !gui.showDetails
			gui.layout()
			return nil
		}

		return event
	})

	return nil
}

func (gui *Gui) nextView() {
	if gui.app.GetFocus() == gui.tools {
		gui.app.SetFocus(gui.conns)
		gui.tools.SetBorderColor(gui.Theme.BorderUnfocus)
		gui.conns.SetBorderColor(gui.Theme.BorderFocus)
		gui.resources.SetBorderColor(gui.Theme.BorderUnfocus)
	} else if gui.app.GetFocus() == gui.conns {
		gui.app.SetFocus(gui.resources)
		gui.tools.SetBorderColor(gui.Theme.BorderUnfocus)
		gui.conns.SetBorderColor(gui.Theme.BorderUnfocus)
		gui.resources.SetBorderColor(gui.Theme.BorderFocus)
	} else {
		gui.app.SetFocus(gui.tools)
		gui.tools.SetBorderColor(gui.Theme.BorderFocus)
		gui.conns.SetBorderColor(gui.Theme.BorderUnfocus)
		gui.resources.SetBorderColor(gui.Theme.BorderUnfocus)
	}
	gui.updateDetails()
}
