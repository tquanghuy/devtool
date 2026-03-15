package checker

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"devtool/internal/config"
	_ "github.com/lib/pq"
)

type PostgresChecker struct{}

func (p *PostgresChecker) Check(ctx context.Context, cfg *config.AppConfig) StatusResult {
	start := time.Now()
	res := StatusResult{
		Tool: ToolTypePostgres,
		Name: "PostgreSQL Database",
	}

	password, _ := config.GetPassword(cfg.Postgres.User)
	
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, password, cfg.Postgres.Database)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		res.IsUp = false
		res.Error = err
		res.Duration = time.Since(start)
		return res
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		res.IsUp = false
		res.Error = err
	} else {
		res.IsUp = true
	}

	res.Duration = time.Since(start)
	return res
}
