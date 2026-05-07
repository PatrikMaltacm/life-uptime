package database

import (
    "database/sql"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib" 
)

func Connect(dsn string) *sql.DB {
    db, err := sql.Open("pgx", dsn)
    if err != nil {
        log.Fatalf("database: failed to open connection: %v", err)
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("database: failed to ping: %v", err)
    }

    log.Println("database: connected to PostgreSQL")

    return db
}