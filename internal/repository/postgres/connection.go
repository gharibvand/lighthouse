package postgres

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func NewConnection(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Performance optimizations
	db.SetMaxOpenConns(25)                        // Maximum number of open connections
	db.SetMaxIdleConns(5)                         // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute)        // 5 minutes - connection lifetime
	db.SetConnMaxIdleTime(10 * time.Minute)       // 10 minutes - idle connection timeout

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
