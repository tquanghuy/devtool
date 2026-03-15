package checker

// ToolKind classifies how many instances of a tool can be managed.
type ToolKind string

const (
	// Singleton tools can only have one managed instance at a time (e.g. Docker, Telepresence).
	Singleton ToolKind = "singleton"
	// PortBound tools can have multiple instances, each bound to a distinct port (e.g. Postgres, MySQL).
	PortBound ToolKind = "portbound"
)

// ToolDefinition is a statically registered tool that devtool knows how to manage.
type ToolDefinition struct {
	// Name is the canonical lowercase identifier used in CLI commands (e.g. "docker", "postgres").
	Name string
	// Kind indicates whether only one or multiple instances are allowed.
	Kind ToolKind
	// DefaultPort is the well-known port for port-bound tools; 0 for singletons.
	DefaultPort int
}
