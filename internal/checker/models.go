package checker

import "time"

type ToolType string

const (
	ToolTypeDocker       ToolType = "Docker"
	ToolTypeTelepresence ToolType = "Telepresence"
	ToolTypePostgres     ToolType = "Postgres"
	ToolTypeMySQL        ToolType = "MySQL"
)





type StatusResult struct {
	Tool     ToolType
	Name     string
	IsUp     bool
	Error    error
	Duration time.Duration
}
