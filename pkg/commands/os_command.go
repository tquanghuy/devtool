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
	// ps -A -o %cpu | awk '{s+=$1} END {print s}'
	out, err := c.RunCommand("ps -A -o %cpu | awk '{s+=$1} END {print s}'")
	if err != nil {
		return "0.0%", nil
	}
	// On Mac, ps %cpu can be > 100% total if multiple cores. 
	// But it's a good relative indicator.
	return out + "%", nil
}

func (c *OSCommand) GetTotalMemUsage() (string, error) {
	// Use vm_stat on macOS
	out, err := c.RunCommand("vm_stat")
	if err != nil {
		return "0 / 0G", nil
	}

	lines := strings.Split(out, "\n")
	var pageSize int64 = 4096 // Default
	if len(lines) > 0 && strings.Contains(lines[0], "page size of") {
		re := regexp.MustCompile(`page size of (\d+) bytes`)
		match := re.FindStringSubmatch(lines[0])
		if len(match) > 1 {
			pageSize, _ = strconv.ParseInt(match[1], 10, 64)
		}
	}

	var active, wired, compressed int64
	for _, line := range lines {
		if strings.HasPrefix(line, "Pages active:") {
			active = c.extractPageCount(line)
		} else if strings.HasPrefix(line, "Pages wired down:") {
			wired = c.extractPageCount(line)
		} else if strings.HasPrefix(line, "Pages occupied by compressor:") {
			compressed = c.extractPageCount(line)
		}
	}

	usedBytes := (active + wired + compressed) * pageSize
	usedG := float64(usedBytes) / 1024 / 1024 / 1024

	// Total memory
	totalRaw, err := c.RunCommand("sysctl -n hw.memsize")
	totalG := "0G"
	if err == nil {
		bytes, _ := strconv.ParseInt(totalRaw, 10, 64)
		totalG = fmt.Sprintf("%dG", bytes/1024/1024/1024)
	}

	return fmt.Sprintf("%.1fG / %s", usedG, totalG), nil
}

func (c *OSCommand) extractPageCount(line string) int64 {
	re := regexp.MustCompile(`(\d+)\.`)
	match := re.FindAllStringSubmatch(line, -1)
	if len(match) > 0 {
		count, _ := strconv.ParseInt(match[len(match)-1][1], 10, 64)
		return count
	}
	return 0
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

// IsPortBusy returns true if the port is already in use
func (c *OSCommand) IsPortBusy(port int) bool {
	return c.DialTCP("localhost", port)
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
