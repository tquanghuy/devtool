package unit

import (
	"testing"

	"devtool/internal/config"
	"devtool/internal/manager"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// noopSelect is a SelectFn that should never be called. Tests that expect
// single-instance removal use this to guard against unexpected prompts.
func noopSelect(ids []string) (int, error) {
	panic("selectFn should not be called for single-instance removal")
}

// selectIndex returns a SelectFn that always picks the given zero-based index.
func selectIndex(i int) manager.SelectFn {
	return func(ids []string) (int, error) {
		return i, nil
	}
}

// TestRemove_ManagedSingletonSucceeds verifies that a managed singleton can be removed.
func TestRemove_ManagedSingletonSucceeds(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("telepresence", 0, false))
	require.NoError(t, manager.RemoveTool("telepresence", false, false, noopSelect, nil))

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	assert.Empty(t, cfg.Instances)
}

// TestRemove_NonManagedToolReturnsError verifies the "not managed" error path.
func TestRemove_NonManagedToolReturnsError(t *testing.T) {
	setTempManagedConfig(t)

	err := manager.RemoveTool("telepresence", false, false, noopSelect, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not currently managed")
}

// TestRemove_MultiInstanceInteractiveRemovesOnlySelected verifies disambiguation
// by injecting a SelectFn that picks the second instance.
func TestRemove_MultiInstanceInteractiveRemovesOnlySelected(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 5432, false))
	require.NoError(t, manager.AddTool("postgres", 5433, false))

	// Select the second instance (idx=1 → "postgres-5433").
	require.NoError(t, manager.RemoveTool("postgres", false, false, selectIndex(1), nil))

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 1)
	assert.Equal(t, "postgres-5432", cfg.Instances[0].Identifier)
}

// TestRemove_NonInteractiveMultipleInstancesFails verifies that −−non-interactive
// returns an error when multiple instances exist and no exact ID is given.
func TestRemove_NonInteractiveMultipleInstancesFails(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 5432, false))
	require.NoError(t, manager.AddTool("postgres", 5433, false))

	err := manager.RemoveTool("postgres", false, true /* non-interactive */, nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Multiple instances")
}

// TestRemove_ByExactIdentifierSucceeds verifies removal using an explicit instance ID.
func TestRemove_ByExactIdentifierSucceeds(t *testing.T) {
	setTempManagedConfig(t)

	require.NoError(t, manager.AddTool("postgres", 5432, false))
	require.NoError(t, manager.AddTool("postgres", 5433, false))

	// Remove by exact identifier — bypasses disambiguation.
	require.NoError(t, manager.RemoveTool("postgres-5432", false, false, noopSelect, nil))

	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.Len(t, cfg.Instances, 1)
	assert.Equal(t, "postgres-5433", cfg.Instances[0].Identifier)
}
