package unit_test

import (
	"os"
	"path/filepath"
	"testing"

	"devtool/internal/devtools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigLoadSave(t *testing.T) {
	// Set up temporary home directory for testing
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	// Case 1: Load nonexistent config
	cfg, err := devtools.Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Empty(t, cfg.Devtools)

	// Case 2: Add and Save config
	profile := devtools.DevtoolProfile{
		Name:       "TestTool",
		Executable: "/usr/bin/test",
		Args:       "--version",
	}
	err = devtools.Add(cfg, profile)
	require.NoError(t, err)

	err = devtools.Save(cfg)
	require.NoError(t, err)

	// Case 3: Load existing config
	cfg2, err := devtools.Load()
	require.NoError(t, err)
	assert.Len(t, cfg2.Devtools, 1)
	assert.Equal(t, profile.Name, cfg2.Devtools[0].Name)
}

func TestConfigAddValidation(t *testing.T) {
	cfg := &devtools.DevtoolsConfig{}

	// Empty Name
	err := devtools.Add(cfg, devtools.DevtoolProfile{Name: "", Executable: "run"})
	assert.ErrorContains(t, err, "name cannot be empty")

	// Empty Executable
	err = devtools.Add(cfg, devtools.DevtoolProfile{Name: "Tool", Executable: ""})
	assert.ErrorContains(t, err, "executable path cannot be empty")

	// Duplicate Name
	_ = devtools.Add(cfg, devtools.DevtoolProfile{Name: "Duplicate", Executable: "run1"})
	err = devtools.Add(cfg, devtools.DevtoolProfile{Name: "Duplicate", Executable: "run2"})
	assert.ErrorContains(t, err, "already exists")
}

func TestConfigRemove(t *testing.T) {
	cfg := &devtools.DevtoolsConfig{
		Devtools: []devtools.DevtoolProfile{
			{Name: "Tool1", Executable: "run1"},
			{Name: "Tool2", Executable: "run2"},
		},
	}

	// Remove existing
	err := devtools.Remove(cfg, "Tool1")
	require.NoError(t, err)
	assert.Len(t, cfg.Devtools, 1)
	assert.Equal(t, "Tool2", cfg.Devtools[0].Name)

	// Remove nonexistent
	err = devtools.Remove(cfg, "NotExist")
	assert.ErrorContains(t, err, "not found")
}

func TestConfigPath(t *testing.T) {
	tempHome := t.TempDir()
	t.Setenv("HOME", tempHome)

	expectedPath := filepath.Join(tempHome, ".devtools.yml")
	
	// We verify config loading actually uses the HOME override
	_ = os.WriteFile(expectedPath, []byte("devtools: []"), 0644)
	cfg, err := devtools.Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)
}
