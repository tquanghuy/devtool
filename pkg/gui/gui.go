package gui

import (
	"github.com/sirupsen/logrus"
	"github.com/rivo/tview"

	"github.com/gdamore/tcell/v2"
	"devtool/pkg/config"
	"devtool/pkg/commands"
)

type State struct {
	Tools []commands.ToolDefinition
}

type Gui struct {
	app     *tview.Application
	pages   *tview.Pages
	tools   *tview.Table
	conns   *tview.Table
	status  *tview.TextView

	Log     *logrus.Entry
	Config  *config.AppConfig
	Managed *config.ManagedConfig
	OS      *commands.OSCommand
	State   State
	Theme   Theme
}

type Theme struct {
	HeaderBg      tcell.Color
	HeaderFg      tcell.Color
	BorderFocus   tcell.Color
	BorderUnfocus tcell.Color
	Running       tcell.Color
	Stopped       tcell.Color
	Connected     tcell.Color
	Disconnected  tcell.Color
}

func defaultTheme() Theme {
	return Theme{
		HeaderBg:      tcell.ColorBlue,
		HeaderFg:      tcell.ColorWhite,
		BorderFocus:   tcell.ColorAqua,
		BorderUnfocus: tcell.ColorGray,
		Running:       tcell.ColorLime,
		Stopped:       tcell.ColorRed,
		Connected:     tcell.ColorAqua,
		Disconnected:  tcell.ColorMaroon,
	}
}

func NewGui(log *logrus.Entry, cfg *config.AppConfig, managed *config.ManagedConfig) (*Gui, error) {
	app := tview.NewApplication()
	
	gui := &Gui{
		app:     app,
		pages:   tview.NewPages(),
		tools:   tview.NewTable(),
		conns:   tview.NewTable(),
		status:  tview.NewTextView(),
		Log:     log,
		Config:  cfg,
		Managed: managed,
		OS:      commands.NewOSCommand(),
	}

	gui.Theme = defaultTheme()

	// Initialize widgets
	gui.tools.SetSelectable(true, false)
	gui.tools.SetBorder(true).
		SetTitle(" Tools ").
		SetBorderColor(gui.Theme.BorderFocus)
	
	gui.conns.SetSelectable(true, false)
	gui.conns.SetBorder(true).
		SetTitle(" Connections ").
		SetBorderColor(gui.Theme.BorderUnfocus)
		
	gui.status.SetDynamicColors(true).SetBackgroundColor(tcell.ColorBlack)

	return gui, nil
}

func (gui *Gui) Run() error {
	gui.layout()

	if err := gui.keybindings(); err != nil {
		return err
	}

	return gui.app.SetRoot(gui.pages, true).Run()
}

func (gui *Gui) Close() {
	gui.app.Stop()
}
