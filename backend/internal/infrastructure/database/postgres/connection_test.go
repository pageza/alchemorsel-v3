package postgres

import "testing"

func TestConnect(t *testing.T) {
	dsn := "postgres://testuser:testpass@localhost/testdb?sslmode=disable"
	db := Connect(dsn)
	if err := db.Ping(); err != nil {
		t.Fatalf("failed to ping db: %v", err)
	}
	db.Close()
}
