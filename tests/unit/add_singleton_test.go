package unit

import (
	"testing"

	"devtool/internal/config"
	"devtool/internal/manager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setTempManagedConfig overrides the home dir so LoadConfig/SaveConfig use a temp dir.
func setTempManagedConfig(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("HOME", dir)
	return dir
}

// TestAddSingleton_Success verifies that adding a singleton for the first time succeeds.
func TestAddSingleton_Success(t *testing.T) {
	setTempManagedConfig(t)

	err := manager.AddTool("docker", 0, false)
	require.NoError(t, err)

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 1)
	assert.Equal(t, "docker", cfg.Instances[0].ToolName)
	assert.Equal(t, "docker", cfg.Instances[0].Identifier)
	assert.False(t, cfg.Instances[0].CreatedAt.IsZero())
}

// TestAddSingleton_DuplicateReturnsError verifies the singleton guard.
func TestAddSingleton_DuplicateReturnsError(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("docker", 0, false))

	err := manager.AddTool("docker", 0, false)
	require.Error(t, err)
	assert.Equal(t, "Tool docker is a singleton and is already managed.", err.Error())
}

// TestAddSingleton_UnsupportedToolReturnsErrorWithList verifies the unsupported-tool path.
func TestAddSingleton_UnsupportedToolReturnsErrorWithList(t *testing.T) {
	setTempManagedConfig(t)

	err := manager.AddTool("notarealtool", 0, false)
	require.Error(t, err)
	assert.Equal(t, "Tool notarealtool is not supported.", err.Error())
}

// TestAddSingleton_TelepresenceAsSingleton validates a second non-docker singleton.
func TestAddSingleton_TelepresenceAsSingleton(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("telepresence", 0, false))

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 1)
	assert.Equal(t, "telepresence", cfg.Instances[0].Identifier)
}

// TestAddSingleton_PersistsAcrossLoads verifies the config is actually written to disk.
func TestAddSingleton_PersistsAcrossLoads(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("docker", 0, false))

	// Reload from disk and confirm the data is there.
	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	assert.Len(t, cfg.Instances, 1)
}

// --- OS-independence helper -----------------------------------------------
// On macOS, HOME is used by os.UserHomeDir internally, so t.Setenv("HOME", …)
// is sufficient. The tests above rely on this.
