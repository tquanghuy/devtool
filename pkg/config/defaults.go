package config

func GetDefaultTools() map[string]ToolDefinition {
	return map[string]ToolDefinition{
		"docker": {
			Name:     "docker",
			CheckCmd: "docker info >/dev/null 2>&1",
			StartCmd: "open -a Docker",
			StopCmd:  "osascript -e 'quit app \"Docker\"'",
		},
		"telepresence": {
			Name:     "telepresence",
			CheckCmd: "telepresence status >/dev/null 2>&1",
			StartCmd: "telepresence connect",
			StopCmd:  "telepresence quit",
		},
		"postgres": {
			Name:        "postgres",
			DefaultPort: 5432,
			CheckCmd:    "pg_isready -h localhost -p %d",
			StartCmd:    "brew services start postgresql",
			StopCmd:     "brew services stop postgresql",
		},
		"mysql": {
			Name:        "mysql",
			DefaultPort: 3306,
			CheckCmd:    "mysqladmin -h localhost -P %d ping",
			StartCmd:    "brew services start mysql",
			StopCmd:     "brew services stop mysql",
		},
		"gcloud": {
			Name:     "gcloud",
			CheckCmd: "gcloud auth list --filter=status:ACTIVE --format='value(account)' | grep -q .",
			StartCmd: "gcloud auth login",
		},
		"cloud-sql-proxy": {
			Name:        "cloud-sql-proxy",
			DefaultPort: 5432,
			CheckCmd:    "lsof -i :%d >/dev/null 2>&1",
			StartCmd:    "cloud-sql-proxy --port %d [INSTANCE_CONNECTION_NAME]",
			StopCmd:     "pkill -f cloud-sql-proxy",
		},
		"pubsub-emulator": {
			Name:        "pubsub-emulator",
			DefaultPort: 8085,
			CheckCmd:    "curl -s localhost:%d >/dev/null 2>&1",
			StartCmd:    "gcloud beta emulators pubsub start --host-port=0.0.0.0:%d",
			StopCmd:     "pkill -f pubsub-emulator",
		},
		"firebase-emulator": {
			Name:        "firebase-emulator",
			DefaultPort: 4000,
			CheckCmd:    "curl -s localhost:%d >/dev/null 2>&1",
			StartCmd:    "firebase emulators:start",
		},
	}
}
