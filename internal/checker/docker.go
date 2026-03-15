package checker

import (
	"context"
	"os/exec"
	"strings"
	"time"

	"devtool/internal/config"
)

type DockerChecker struct{}

func (d *DockerChecker) Check(ctx context.Context, cfg *config.AppConfig) StatusResult {
	start := time.Now()
	res := StatusResult{
		Tool: ToolTypeDocker,
		Name: "Docker Daemon",
	}

	cmd := exec.CommandContext(ctx, "docker", "info")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		res.IsUp = false
		res.Error = err
	} else {
		res.IsUp = ParseDockerInfo(string(output))
	}

	res.Duration = time.Since(start)
	return res
}

func ParseDockerInfo(output string) bool {
	// If "Cannot connect to the Docker daemon" is in output, it means daemon is not running.
	return !strings.Contains(output, "Cannot connect to the Docker daemon")
}
