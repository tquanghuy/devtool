package commands

import (
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ResourceStat struct {
	Name string
	Type string
	CPU  string
	MEM  string
}

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

// KillPID kills a process by PID
func (c *OSCommand) KillPID(pid string) error {
	_, err := c.RunCommand(fmt.Sprintf("kill -9 %s", pid))
	return err
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

func (c *OSCommand) GetDockerStats() ([]ResourceStat, error) {
	out, err := c.RunCommand("docker stats --no-stream --format '{{.Name}}|{{.CPUPerc}}|{{.MemUsage}}'")
	if err != nil {
		return nil, err
	}

	var stats []ResourceStat
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) >= 3 {
			stats = append(stats, ResourceStat{
				Name: parts[0],
				Type: "Docker",
				CPU:  parts[1],
				MEM:  parts[2],
			})
		}
	}
	return stats, nil
}

func (c *OSCommand) GetTotalCPUUsage() (string, error) {
	// top -l 1 | grep "CPU usage" | awk '{print $3 + $5}'
	out, err := c.RunCommand("top -l 1 | grep 'CPU usage' | awk '{print $3 + $5}'")
	if err != nil {
		return "0.0%", nil
	}
	return out + "%", nil
}

func (c *OSCommand) GetTotalMemUsage() (string, error) {
	// Used memory
	used, err := c.RunCommand("top -l 1 | grep 'PhysMem' | awk '{print $2}'")
	if err != nil {
		used = "0G"
	}
	
	// Total memory
	totalRaw, err := c.RunCommand("sysctl -n hw.memsize")
	total := "0G"
	if err == nil {
		bytes, _ := strconv.ParseInt(totalRaw, 10, 64)
		total = fmt.Sprintf("%dG", bytes/1024/1024/1024)
	}
	
	return fmt.Sprintf("%s / %s", used, total), nil
}

func (c *OSCommand) GetTopProcesses(limit int) ([]ResourceStat, error) {
	cmd := "ps -Ao pcpu,pmem,pid,command -ww -r"
	if limit > 0 {
		cmd = fmt.Sprintf("%s | head -n %d", cmd, limit+1)
	}

	out, err := c.RunCommand(cmd)
	if err != nil {
		return nil, err
	}

	var stats []ResourceStat
	lines := strings.Split(out, "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	re := regexp.MustCompile(`\s+`)
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Format: %CPU %MEM PID COMMAND
		parts := re.Split(line, 4)
		if len(parts) >= 4 {
			stats = append(stats, ResourceStat{
				Name: parts[3], // Full command
				Type: parts[2], // PID
				CPU:  parts[0] + "%",
				MEM:  parts[1] + "%",
			})
		}
	}
	return stats, nil
}

func (c *OSCommand) GetProcessStats(processName string) (*ResourceStat, error) {
	// Find all PIDs for the process name
	out, err := c.RunCommand(fmt.Sprintf("pgrep -f %s", processName))
	if err != nil || out == "" {
		return nil, nil // Not running
	}

	pids := strings.ReplaceAll(out, "\n", ",")
	
	// Get aggregated CPU and MEM for these PIDs
	// ps -p 123,456 -o %cpu,%mem --no-headers (Linux/Mac might differ slightly)
	// On Mac, -o usually works.
	statsOut, err := c.RunCommand(fmt.Sprintf("ps -p %s -o %%cpu,%%mem", pids))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(statsOut, "\n")
	if len(lines) <= 1 {
		return nil, nil
	}

	var totalCPU float64
	var totalMEM float64

	re := regexp.MustCompile(`\s+`)
	for _, line := range lines[1:] { // Skip header
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := re.Split(line, -1)
		if len(parts) >= 2 {
			cpu, _ := strconv.ParseFloat(parts[0], 64)
			mem, _ := strconv.ParseFloat(parts[1], 64)
			totalCPU += cpu
			totalMEM += mem
		}
	}

	return &ResourceStat{
		Name: processName,
		Type: "Local",
		CPU:  fmt.Sprintf("%.1f%%", totalCPU),
		MEM:  fmt.Sprintf("%.1f%%", totalMEM),
	}, nil
}

// GetFreePort finds the next available port starting from startPort
func (c *OSCommand) GetFreePort(startPort int) int {
	port := startPort
	for {
		address := fmt.Sprintf(":%d", port)
		l, err := net.Listen("tcp", address)
		if err == nil {
			l.Close()
			return port
		}
		port++
		if port > 65535 {
			return 0 // Failed to find a port
		}
	}
}

// FormatCommand replaces %d placeholders with the actual port
func (c *OSCommand) FormatCommand(cmd string, port int) string {
	if !strings.Contains(cmd, "%d") {
		return cmd
	}
	return fmt.Sprintf(cmd, port)
}
