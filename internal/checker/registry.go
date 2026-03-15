package checker

import "strings"

// registeredTools is the authoritative list of tools devtool can manage.
var registeredTools = []ToolDefinition{
	{Name: "docker", Kind: Singleton, DefaultPort: 0},
	{Name: "telepresence", Kind: Singleton, DefaultPort: 0},
	{Name: "postgres", Kind: PortBound, DefaultPort: 5432},
	{Name: "mysql", Kind: PortBound, DefaultPort: 3306},
}

// LookupTool returns the ToolDefinition for the given name (case-insensitive).
// The second return value is false if the tool is not registered.
func LookupTool(name string) (*ToolDefinition, bool) {
	lower := strings.ToLower(name)
	for i := range registeredTools {
		if registeredTools[i].Name == lower {
			return &registeredTools[i], true
		}
	}
	return nil, false
}

// SupportedToolNames returns the sorted list of all registered tool names.
func SupportedToolNames() []string {
	names := make([]string, len(registeredTools))
	for i, t := range registeredTools {
		names[i] = t.Name
	}
	return names
}
