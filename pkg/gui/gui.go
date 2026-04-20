package gui

import (
	"fmt"
	"sync"
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
		for range ticker.C {
			var (
				cpu, mem     string
				resources    []commands.ResourceStat
				toolStatuses = make(map[string]string)
				connStatuses = make(map[string]string)
				allConns     = []ConnectionSpec{}
			)

			var wg sync.WaitGroup

			// 1. System Stats
			wg.Add(1)
			go func() {
				defer wg.Done()
				cpu, _ = gui.OS.GetTotalCPUUsage()
				mem, _ = gui.OS.GetTotalMemUsage()
			}()

			// 2. Top Processes
			wg.Add(1)
			go func() {
				defer wg.Done()
				resources, _ = gui.OS.GetTopProcesses(50) // Reduced to 50 for performance
			}()

			// 3. Tool Statuses
			wg.Add(1)
			go func() {
				defer wg.Done()
				// We can even parallelize per-tool checks if we want
				type result struct {
					id     string
					status string
				}
				resChan := make(chan result, len(gui.State.Tools))
				var toolWg sync.WaitGroup
				for _, tool := range gui.State.Tools {
					toolWg.Add(1)
					go func(t ToolInstance) {
						defer toolWg.Done()
						status := "STOPPED"
						checkCmd := gui.OS.FormatCommand(t.Definition.CheckCmd, t.Instance.Port)
						if gui.OS.CheckToolStatus(checkCmd) {
							status = "RUNNING"
						}
						resChan <- result{t.Instance.Identifier, status}
					}(tool)
				}
				toolWg.Wait()
				close(resChan)
				for r := range resChan {
					toolStatuses[r.id] = r.status
				}
			}()

			// 4. Connection Statuses
			wg.Add(1)
			go func() {
				defer wg.Done()
				for name, conn := range gui.Config.PostgresConns {
					allConns = append(allConns, ConnectionSpec{Name: name, Host: conn.Host, Port: conn.Port})
				}
				for name, conn := range gui.Config.MySQLConns {
					allConns = append(allConns, ConnectionSpec{Name: name, Host: conn.Host, Port: conn.Port})
				}

				type result struct {
					name   string
					status string
				}
				resChan := make(chan result, len(allConns))
				var connWg sync.WaitGroup
				for _, conn := range allConns {
					connWg.Add(1)
					go func(c ConnectionSpec) {
						defer connWg.Done()
						status := gui.checkConnection("", c.Host, c.Port)
						resChan <- result{c.Name, status}
					}(conn)
				}
				connWg.Wait()
				close(resChan)
				for r := range resChan {
					connStatuses[r.name] = r.status
				}
			}()

			wg.Wait()

			gui.app.QueueUpdateDraw(func() {
				gui.State.TotalCPU = cpu
				gui.State.TotalMem = mem
				gui.State.Resources = resources
				gui.State.ToolStatuses = toolStatuses
				gui.State.ConnStatuses = connStatuses
				gui.State.Connections = allConns

				gui.renderResources()
				gui.renderTools()
				
				// Also update details to reflect new status if something changed
				gui.updateDetails()
			})
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
