package postgres

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"

	"alchemorsel/backend/internal/config"
	"github.com/fergusstrange/embedded-postgres"
)

func TestConnectAndMigrate(t *testing.T) {
	ep := embeddedpostgres.NewDatabase()
	if err := ep.Start(); err != nil {
		t.Skipf("embedded postgres unavailable: %v", err)
	}
	defer ep.Stop()

	cfg := config.DatabaseConfig{
		Host:         "localhost",
		Port:         5432,
		User:         "postgres",
		Password:     "postgres",
		Name:         "postgres",
		SSLMode:      "disable",
		MaxOpenConns: 7,
		MaxIdleConns: 3,
	}

	migrationsPath := filepath.Join("internal", "infrastructure", "database", "postgres", "migrations")

	db, err := Connect(cfg, migrationsPath)
	if err != nil {
		t.Fatalf("connect failed: %v", err)
	}
	defer db.Close()

	if _, err := db.Query("SELECT 1 FROM test_table LIMIT 1"); err != nil {
		t.Fatalf("migration not applied: %v", err)
	}

	if db.Stats().MaxOpenConnections != cfg.MaxOpenConns {
		t.Fatalf("expected max open %d, got %d", cfg.MaxOpenConns, db.Stats().MaxOpenConnections)
	}

	conns := []*sql.Conn{}
	for i := 0; i < cfg.MaxIdleConns+2; i++ {
		c, err := db.Conn(context.Background())
		if err != nil {
			t.Fatalf("getting conn: %v", err)
		}
		conns = append(conns, c)
	}
	for _, c := range conns {
		c.Close()
	}
	if idle := db.Stats().Idle; idle > cfg.MaxIdleConns {
		t.Fatalf("idle connections %d exceed max idle %d", idle, cfg.MaxIdleConns)
	}
}
