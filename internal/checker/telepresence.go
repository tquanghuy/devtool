package checker

import (
	"context"
	"errors"
	"os/exec"
	"strings"
	"time"

	"devtool/internal/config"
)

type TelepresenceChecker struct{}

func (t *TelepresenceChecker) Check(ctx context.Context, cfg *config.AppConfig) StatusResult {
	start := time.Now()
	res := StatusResult{
		Tool: ToolTypeTelepresence,
		Name: "Telepresence",
	}

	cmd := exec.CommandContext(ctx, "telepresence", "status")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		res.IsUp = false
		res.Error = err
	} else {
		res.IsUp = ParseTelepresenceStatus(string(output))
		if !res.IsUp {
			res.Error = errors.New("daemon not running")
		}
	}

	res.Duration = time.Since(start)
	return res
}

func ParseTelepresenceStatus(output string) bool {
	// Telepresence status will output "Not running" if the daemon is stopped
	return !strings.Contains(output, "Not running")
}
