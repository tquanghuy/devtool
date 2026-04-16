package commands

import (
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

// OSCommand is a wrapper around shell commands
type OSCommand struct {
	// Add logrus reference later if needed
}

func NewOSCommand() *OSCommand {
	return &OSCommand{}
}

// RunCommand executes a command and returns the output
func (c *OSCommand) RunCommand(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// IsRunning checks if a process is running
func (c *OSCommand) IsRunning(processName string) bool {
	_, err := exec.Command("pgrep", "-i", processName).Output()
	return err == nil
}

// KillProcess kills a process by name
func (c *OSCommand) KillProcess(processName string) error {
	return exec.Command("pkill", "-i", processName).Run()
}
// CheckToolStatus runs the check command and returns true if it exits with 0
func (c *OSCommand) CheckToolStatus(checkCmd string) bool {
	if checkCmd == "" {
		return false
	}
	cmd := exec.Command("sh", "-c", checkCmd)
	err := cmd.Run()
	return err == nil
}

// DialTCP checks if a TCP port is open
func (c *OSCommand) DialTCP(host string, port int) bool {
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
