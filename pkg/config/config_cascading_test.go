package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_CascadingLoad(t *testing.T) {
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

	// 2. Setup local config in temp project dir
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
  global-tool:
    name: "overridden-global"
    kind: "singleton"
    check_cmd: "echo overridden"
`
	err = os.WriteFile(".devtool.yml", []byte(localYAML), 0644)
	require.NoError(t, err)

	// 3. Load cascaded config
	cfg, err := Load()
	require.NoError(t, err)

	// Verify defaults are present
	assert.Contains(t, cfg.Tools, "docker")
	
	// Verify global tool is present but overridden
	assert.Contains(t, cfg.Tools, "global-tool")
	assert.Equal(t, "echo overridden", cfg.Tools["global-tool"].CheckCmd)
	
	// Verify local-only tool is present
	assert.Contains(t, cfg.Tools, "local-tool")
	assert.Equal(t, "echo local", cfg.Tools["local-tool"].CheckCmd)
}
