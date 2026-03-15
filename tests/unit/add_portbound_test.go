package unit

import (
	"testing"

	"devtool/internal/config"
	"devtool/internal/manager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAddPortbound_FirstAddSucceeds verifies the first postgres add uses default port 5432.
func TestAddPortbound_FirstAddSucceeds(t *testing.T) {
	setTempManagedConfig(t)

	err := manager.AddTool("postgres", 0, false)
	require.NoError(t, err)

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 1)
	assert.Equal(t, "postgres-5432", cfg.Instances[0].Identifier)
}

// TestAddPortbound_SecondAddWithoutPortFails verifies the conflict error when
// adding the same port-bound tool twice without specifying −−port.
func TestAddPortbound_SecondAddWithoutPortFails(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 0, false))

	err := manager.AddTool("postgres", 0, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Port conflict for postgres. Please specify a different port or run interactively.")
}

// TestAddPortbound_ExplicitFreePorts verifies that two different explicit ports succeed.
func TestAddPortbound_ExplicitFreePorts(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 5432, false))
	require.NoError(t, manager.AddTool("postgres", 5433, false))

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 2)
}

// TestAddPortbound_ExplicitPortConflict verifies that using the same explicit port twice fails.
func TestAddPortbound_ExplicitPortConflict(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 5435, false))

	err := manager.AddTool("postgres", 5435, false)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Port conflict for postgres. Please specify a different port or run interactively.")
}

// TestAddPortbound_NonInteractiveOnConflict verifies −−non-interactive fails without a prompt.
// (The current implementation always returns error on conflict regardless of the flag,
// so this test verifies the error is returned when nonInteractive is true.)
func TestAddPortbound_NonInteractiveOnConflict(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 0, false))

	err := manager.AddTool("postgres", 0, true /* non-interactive */)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Port conflict for postgres. Please specify a different port or run interactively.")
}

// TestAddPortbound_MySQL verifies a second port-bound tool type works correctly.
func TestAddPortbound_MySQL(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("mysql", 0, false))

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 1)
	assert.Equal(t, "mysql-3306", cfg.Instances[0].Identifier)
}
