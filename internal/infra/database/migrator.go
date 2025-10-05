package database

import (
	"embed"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

//go:embed migrations/**/*.sql
var migrationsFS embed.FS

type Migrator struct {
	db     *sqlx.DB
	driver string
}

// TODO: create a specific logger
func NewMigrator(db *sqlx.DB, driver string) *Migrator {
	m := &Migrator{db: db, driver: driver}
	m.ensureMigrationsTable()
	return m
}

func (m *Migrator) ensureMigrationsTable() {
	schema := `
		CREATE TABLE IF NOT EXISTS migrations (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			batch INTEGER NOT NULL,
			applied_at TIMESTAMP NOT NULL
		);
	`
	if _, err := m.db.Exec(schema); err != nil {
		panic(fmt.Errorf("failed to create migrations: %w", err))
	}
}

type migration struct {
	ID   string
	Name string
	Up   string
	Down string
}

func (m *Migrator) Up() error {
	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	allMigrations, err := m.loadMigrations()
	if err != nil {
		return err
	}

	// get last batch
	var lastBatch int
	err = m.db.Get(&lastBatch, "SELECT COALESCE(MAX(batch),0) FROM migrations;")
	if err != nil {
		return err
	}

	currentBatch := lastBatch + 1

	for _, mig := range allMigrations {
		if _, ok := applied[mig.ID]; ok {
			continue // already applied
		}

		if strings.TrimSpace(mig.Up) == "" {
			return fmt.Errorf("no up migration found for %s", mig.Name)
		}

		if _, err := m.db.Exec(strings.TrimSpace(mig.Up)); err != nil {
			return err
		}

		_, err = m.db.NamedExec(
			`INSERT INTO migrations (id, name, batch, applied_at) VALUES (:id, :name, :batch, :applied_at);`,
			map[string]interface{}{
				"id":         mig.ID,
				"name":       mig.Name,
				"batch":      currentBatch,
				"applied_at": time.Now().UTC(),
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) Down(n int) error {
	rows := []struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		Batch int    `db:"batch"`
	}{}

	nstmt, err := m.db.PrepareNamed(`SELECT id, name, batch FROM migrations ORDER BY batch DESC, applied_at DESC LIMIT :limit;`)
	if err != nil {
		return err
	}
	defer nstmt.Close()

	if err := nstmt.Select(&rows, map[string]interface{}{
		"limit": n,
	}); err != nil {
		return err
	}

	migs, err := m.loadMigrations()
	if err != nil {
		return err
	}

	migMap := map[string]migration{}
	for _, mig := range migs {
		migMap[mig.ID] = mig
	}

	for _, row := range rows {
		mig, ok := migMap[row.ID]
		if !ok || strings.TrimSpace(mig.Down) == "" {
			return fmt.Errorf("no down migration found for %s", row.Name)
		}

		if _, err := m.db.Exec(strings.TrimSpace(mig.Down)); err != nil {
			return err
		}

		_, err := m.db.NamedExec("DELETE FROM migrations WHERE id = :id", map[string]interface{}{
			"id": row.ID,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Migrator) Reset() error {
	var count int
	err := m.db.Get(&count, "SELECT COUNT(*) FROM migrations;")
	if err != nil {
		fmt.Println("Error during migration count retrieval:", err)
		return err
	}

	if err := m.Down(count); err != nil {
		fmt.Println("Error during rollback:", err)
		return err
	}

	return m.Up()
}

func (m *Migrator) Rollback() error {
	var lastBatch int
	err := m.db.Get(&lastBatch, "SELECT COALESCE(MAX(batch),0) FROM migrations;")
	if err != nil {
		return err
	}

	if lastBatch == 0 {
		return nil
	}

	var count int
	nstmt, err := m.db.PrepareNamed("SELECT COUNT(*) FROM migrations WHERE batch = :batch;")
	if err != nil {
		return err
	}
	defer nstmt.Close()

	err = nstmt.Get(&count, map[string]interface{}{
		"batch": lastBatch,
	})
	if err != nil {
		return err
	}

	return m.Down(count)
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	rows := []struct {
		ID string `db:"id"`
	}{}

	err := m.db.Select(&rows, "SELECT id FROM migrations;")
	if err != nil {
		return nil, err
	}

	applied := map[string]bool{}
	for _, r := range rows {
		applied[r.ID] = true
	}
	return applied, nil
}

func (m *Migrator) loadMigrations() ([]migration, error) {
	entries, err := migrationsFS.ReadDir("migrations/" + m.driver)
	if err != nil {
		return nil, err
	}

	migs := []migration{}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		fileName := e.Name()
		parts := strings.SplitN(fileName, ".", 3)
		if len(parts) < 3 {
			continue
		}

		id := parts[0]
		name := parts[1]

		dir := "migrations/" + m.driver + "/" + fileName
		sqlBytes, err := migrationsFS.ReadFile(dir)
		if err != nil {
			return nil, err
		}

		content := string(sqlBytes)

		var up, down string
		if strings.Contains(content, "-- DOWN") {
			parts := strings.SplitN(content, "-- DOWN", 2)
			up = strings.TrimSpace(parts[0])
			if len(parts) > 1 {
				down = strings.TrimSpace(parts[1])
			}
		} else {
			up = content
		}

		migs = append(migs, migration{
			ID:   id,
			Name: name,
			Up:   up,
			Down: down,
		})
	}

	sort.Slice(migs, func(i, j int) bool { return migs[i].ID < migs[j].ID })
	return migs, nil
}
