package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"           // Postgres driver
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

func New(config *config.DatabaseConfig) *sqlx.DB {
	db, err := sqlx.Open(config.CodeDriver(), config.GetDSN())
	if err != nil {
		panic(fmt.Errorf("failed to open db: %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping db: %w", err))
	}

	return db
}
