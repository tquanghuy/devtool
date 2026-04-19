package gui

import (
	"fmt"
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

type ToolInstance struct {
	Definition config.ToolDefinition
	Instance   config.ManagedInstance
}

type State struct {
	Tools        []ToolInstance
	ToolStatuses map[string]string
	Connections  []ConnectionSpec
	ConnStatuses map[string]string
	Resources    []commands.ResourceStat
	TotalCPU     string
	TotalMem     string
}

type Gui struct {
	app       *tview.Application
	pages     *tview.Pages
	tools     *tview.Table
	resources *tview.Table
	details   *tview.TextView
	status    *tview.TextView

	showDetails bool
	Log         *logrus.Entry
	Config      *config.AppConfig
	OS          *commands.OSCommand
	State       State
	Theme       Theme
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

func NewGui(log *logrus.Entry, cfg *config.AppConfig) (*Gui, error) {
	gui := &Gui{
		app:         tview.NewApplication(),
		pages:       tview.NewPages(),
		tools:       tview.NewTable(),
		resources:   tview.NewTable(),
		details:     tview.NewTextView(),
		status:      tview.NewTextView(),
		Log:         log,
		Config:      cfg,
		OS:          commands.NewOSCommand(),
		showDetails: true,
		Theme:       defaultTheme(),
		State: State{
			ToolStatuses: make(map[string]string),
			ConnStatuses: make(map[string]string),
		},
	}

	// Initialize widgets
	gui.tools.SetSelectable(true, false).
		SetBorder(true).
		SetTitle(" Tools ").
		SetBorderColor(gui.Theme.BorderFocus)
	gui.tools.SetSelectionChangedFunc(func(row, column int) {
		gui.updateDetails()
	})

	gui.resources.SetSelectable(true, false).
		SetBorder(true).
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

	// Background polling for resources and statuses
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-ticker.C:
				// 1. Get System Stats
				cpu, _ := gui.OS.GetTotalCPUUsage()
				mem, _ := gui.OS.GetTotalMemUsage()

				// 2. Get Top Processes
				topProcs, _ := gui.OS.GetTopProcesses(100)

				// 3. Poll Tool Statuses
				toolStatuses := make(map[string]string)
				for _, tool := range gui.State.Tools {
					status := "STOPPED"
					checkCmd := gui.OS.FormatCommand(tool.Definition.CheckCmd, tool.Instance.Port)
					if gui.OS.CheckToolStatus(checkCmd) {
						status = "RUNNING"
					}
					toolStatuses[tool.Instance.Identifier] = status
				}

				// 4. Poll Connection Statuses (all)
				connStatuses := make(map[string]string)
				allConns := []ConnectionSpec{}
				for name, conn := range gui.Config.PostgresConns {
					allConns = append(allConns, ConnectionSpec{Name: name, Host: conn.Host, Port: conn.Port})
				}
				for name, conn := range gui.Config.MySQLConns {
					allConns = append(allConns, ConnectionSpec{Name: name, Host: conn.Host, Port: conn.Port})
				}

				for _, conn := range allConns {
					status := gui.checkConnection("", conn.Host, conn.Port)
					connStatuses[conn.Name] = status
				}

				gui.app.QueueUpdateDraw(func() {
					gui.State.TotalCPU = cpu
					gui.State.TotalMem = mem
					gui.State.Resources = topProcs
					gui.State.ToolStatuses = toolStatuses
					gui.State.ConnStatuses = connStatuses
					gui.State.Connections = allConns
					
					gui.renderResources()
					gui.renderTools()
				})
			}
		}
	}()

	return gui.app.SetRoot(gui.pages, true).Run()
}

func (gui *Gui) Close() {
	gui.app.Stop()
}

func (gui *Gui) handleQuit() {
	modal := tview.NewModal().
		SetText("Are you sure you want to quit?").
		AddButtons([]string{"Quit", "Cancel"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == "Quit" {
				gui.app.Stop()
			}
			gui.pages.RemovePage("quit")
		})

	gui.pages.AddPage("quit", modal, true, true)
}

func (gui *Gui) checkConnection(toolName string, host string, port int) string {
	if toolName != "" {
		if t, ok := config.GetDefaultTools()[toolName]; ok {
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
