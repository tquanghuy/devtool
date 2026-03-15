package config

import "time"

// ManagedInstance represents a single tool instance in the managed list.
type ManagedInstance struct {
	// ToolName is the canonical tool name (matches ToolDefinition.Name).
	ToolName string `json:"tool_name"`
	// Identifier is the unique key for this instance.
	// For singletons it equals ToolName; for port-bound tools it is "<ToolName>-<port>".
	Identifier string `json:"identifier"`
	// CreatedAt records when this instance was added.
	CreatedAt time.Time `json:"created_at"`
}
