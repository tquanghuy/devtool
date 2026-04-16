package commands

import (
	"devtool/pkg/config"
)

func GetSupportedToolNames() []string {
	defaults := config.GetDefaultTools()
	names := make([]string, 0, len(defaults))
	for name := range defaults {
		names = append(names, name)
	}
	return names
}

func LookupTool(name string) (config.ToolDefinition, bool) {
	def, ok := config.GetDefaultTools()[name]
	return def, ok
}


