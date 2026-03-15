package checker

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"devtool/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQLChecker struct{}

func (m *MySQLChecker) Check(ctx context.Context, cfg *config.AppConfig) StatusResult {
	start := time.Now()
	res := StatusResult{
		Tool: ToolTypeMySQL,
		Name: "MySQL Database",
	}

	password, _ := config.GetPassword(cfg.MySQL.User)
	
	// Format: user:password@tcp(host:port)/dbname
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.MySQL.User, password, cfg.MySQL.Host, cfg.MySQL.Port, cfg.MySQL.Database)

	db, err := sql.Open("mysql", connStr)
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
