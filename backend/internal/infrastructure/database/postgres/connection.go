package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type DB struct {
	*sql.DB
}

func Connect(connStr string) *DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return &DB{db}
}
