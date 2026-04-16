package gui

import (
	"time"

	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"

	"devtool/pkg/commands"
	"devtool/pkg/config"

	"github.com/gdamore/tcell/v2"
)

type ConnectionSpec struct {
	Name string
	Host string
	Port int
}

type State struct {
	Tools       []config.ToolDefinition
	Connections []ConnectionSpec
	Resources   []commands.ResourceStat
	TotalCPU    string
	TotalMem    string
}

type Gui struct {
	app       *tview.Application
	pages     *tview.Pages
	tools     *tview.Table
	conns     *tview.Table
	resources *tview.Table
	details   *tview.TextView
	status    *tview.TextView

	showDetails bool
	Log         *logrus.Entry
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
		app:       app,
		pages:     tview.NewPages(),
		tools:     tview.NewTable(),
		conns:     tview.NewTable(),
		resources: tview.NewTable(),
		details:   tview.NewTextView(),
		status:    tview.NewTextView(),
		Log:       log,
		Config:    cfg,
		Managed:   managed,
		OS:        commands.NewOSCommand(),
		showDetails: true,
	}

	gui.Theme = defaultTheme()

	// Initialize widgets
	gui.tools.SetSelectable(true, false)
	gui.tools.SetBorder(true).
		SetTitle(" Tools ").
		SetBorderColor(gui.Theme.BorderFocus)
	gui.tools.SetSelectionChangedFunc(func(row, column int) {
		gui.updateDetails()
	})

	gui.conns.SetSelectable(true, false)
	gui.conns.SetBorder(true).
		SetTitle(" Connections ").
		SetBorderColor(gui.Theme.BorderUnfocus)
	gui.conns.SetSelectionChangedFunc(func(row, column int) {
		gui.updateDetails()
	})

	gui.resources.SetSelectable(true, false)
	gui.resources.SetBorder(true).
		SetTitle(" Resource Monitoring ").
		SetBorderColor(gui.Theme.BorderUnfocus)
	gui.resources.SetSelectionChangedFunc(func(row, column int) {
		gui.updateDetails()
	})

	gui.details.SetDynamicColors(true).
		SetBorder(true).
		SetTitle(" Details ").
		SetBorderColor(gui.Theme.BorderUnfocus)

	gui.status.SetDynamicColors(true).SetBackgroundColor(tcell.ColorBlack)

	return gui, nil
}

func (gui *Gui) Run() error {
	gui.layout()

	if err := gui.keybindings(); err != nil {
		return err
	}

	// Background polling for resources
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				// 1. Get System Stats
				cpu, _ := gui.OS.GetTotalCPUUsage()
				mem, _ := gui.OS.GetTotalMemUsage()

				// 2. Get Top Processes (Increased limit for scrolling)
				topProcs, _ := gui.OS.GetTopProcesses(100)

				gui.app.QueueUpdateDraw(func() {
					gui.State.TotalCPU = cpu
					gui.State.TotalMem = mem
					gui.State.Resources = topProcs
					gui.renderResources()
				})
			}
		}
	}()

	return gui.app.SetRoot(gui.pages, true).Run()
}

func (gui *Gui) Close() {
	gui.app.Stop()
}
