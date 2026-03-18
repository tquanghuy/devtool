package devtools

// DevtoolProfile represents a user-configured developer tool entry.
type DevtoolProfile struct {
	Name       string `yaml:"name"`
	Executable string `yaml:"executable"`
	Args       string `yaml:"args"`
}

// DevtoolsConfig is the top-level configuration for all user devtools.
type DevtoolsConfig struct {
	Devtools []DevtoolProfile `yaml:"devtools"`
}
