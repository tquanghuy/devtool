package commands

type ToolKind int

const (
	Singleton ToolKind = iota
	PortBound
)

type ToolDefinition struct {
	Name        string
	Kind        ToolKind
	DefaultPort int
	CheckCmd    string
	StartCmd    string
	StopCmd     string
}

var supportedTools = map[string]ToolDefinition{
	"docker": {
		Name:     "docker",
		Kind:     Singleton,
		CheckCmd: "docker info >/dev/null 2>&1",
		StartCmd: "open -a Docker",
		StopCmd:  "osascript -e 'quit app \"Docker\"'",
	},
	"telepresence": {
		Name:     "telepresence",
		Kind:     Singleton,
		CheckCmd: "telepresence status >/dev/null 2>&1",
		StartCmd: "telepresence connect",
		StopCmd:  "telepresence quit",
	},
	"postgres": {
		Name:        "postgres",
		Kind:        PortBound,
		DefaultPort: 5432,
		CheckCmd:    "pg_isready -h localhost -p %d",
		StartCmd:    "brew services start postgresql",
		StopCmd:     "brew services stop postgresql",
	},
	"mysql": {
		Name:        "mysql",
		Kind:        PortBound,
		DefaultPort: 3306,
		CheckCmd:    "mysqladmin -h localhost -P %d ping",
		StartCmd:    "brew services start mysql",
		StopCmd:     "brew services stop mysql",
	},
}

func GetSupportedToolNames() []string {
	names := make([]string, 0, len(supportedTools))
	for name := range supportedTools {
		names = append(names, name)
	}
	return names
}

func LookupTool(name string) (ToolDefinition, bool) {
	def, ok := supportedTools[name]
	return def, ok
}
