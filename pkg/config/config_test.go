package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_GlobalLoad(t *testing.T) {
	// 1. Setup temporary home for global config
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)
	
	globalConfigPath := filepath.Join(tempHome, ".devtool.yml")
	globalYAML := `
tools:
  global-tool:
    name: "global-tool"
    kind: "singleton"
    check_cmd: "echo global"
`
	err := os.WriteFile(globalConfigPath, []byte(globalYAML), 0644)
	require.NoError(t, err)

	// 2. Setup local config in temp project dir (which should be ignored)
	tempProject := t.TempDir()
	originalWD, _ := os.Getwd()
	os.Chdir(tempProject)
	defer os.Chdir(originalWD)

	localYAML := `
tools:
  local-tool:
    name: "local-tool"
    kind: "singleton"
    check_cmd: "echo local"
`
	err = os.WriteFile(".devtool.yml", []byte(localYAML), 0644)
	require.NoError(t, err)

	// 3. Load config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify defaults are present
	assert.Contains(t, cfg.Tools, "docker")
	
	// Verify global tool is present
	assert.Contains(t, cfg.Tools, "global-tool")
	assert.Equal(t, "echo global", cfg.Tools["global-tool"].CheckCmd)
	
	// Verify local-only tool is NOT present (ignored)
	assert.NotContains(t, cfg.Tools, "local-tool")
}
