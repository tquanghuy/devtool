package tui

import (
	"bytes"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TelepresenceStatus represents the connection status of the Telepresence CLI
type TelepresenceStatus int

const (
	StatusUnknown TelepresenceStatus = iota
	StatusNotInstalled
	StatusDisconnected
	StatusConnected
)

// CheckTelepresenceInstalled checks if the telepresence binary is in PATH
func CheckTelepresenceInstalled() bool {
	_, err := exec.LookPath("telepresence")
	return err == nil
}

// TelepresenceStatusMsg is a tea.Msg containing the current status
type TelepresenceStatusMsg struct {
	Status TelepresenceStatus
	Err    error
}

// CheckTelepresenceStatusCmd is a tea.Cmd that runs `telepresence status`
// and parses the output to determine connection status.
func CheckTelepresenceStatusCmd() tea.Cmd {
	return func() tea.Msg {
		if !CheckTelepresenceInstalled() {
			return TelepresenceStatusMsg{Status: StatusNotInstalled, Err: nil}
		}

		cmd := exec.Command("telepresence", "status")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		err := cmd.Run()
		output := out.String()

		if err != nil {
			// Sometimes telepresence returns error code if daemon is not running
			if strings.Contains(output, "Not running") || strings.Contains(output, "telepresence daemon is not running") || strings.Contains(output, "not running") {
				return TelepresenceStatusMsg{Status: StatusDisconnected, Err: nil}
			}
			// Connection issues or other errors
			return TelepresenceStatusMsg{Status: StatusUnknown, Err: err}
		}

		// Parsing the output
		if strings.Contains(output, "Root Daemon: Running") || strings.Contains(output, "User Daemon: Running") || strings.Contains(output, "telepresence is running") || strings.Contains(output, "Connected") {
			return TelepresenceStatusMsg{Status: StatusConnected, Err: nil}
		}
		
		if strings.Contains(output, "Not running") {
			return TelepresenceStatusMsg{Status: StatusDisconnected, Err: nil}
		}

		// Fallback
		return TelepresenceStatusMsg{Status: StatusUnknown, Err: nil}
	}
}

type TelepresenceActionMsg struct {
	Action string
	Err    error
}

func ConnectTelepresenceCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("telepresence", "connect")
		err := cmd.Run()
		return TelepresenceActionMsg{Action: "connect", Err: err}
	}
}

func DisconnectTelepresenceCmd() tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("telepresence", "quit")
		err := cmd.Run()
		return TelepresenceActionMsg{Action: "disconnect", Err: err}
	}
}
