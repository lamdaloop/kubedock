package api

import (
    "fmt"
    "log"
    "os"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

var DB *sqlx.DB

func InitDB() {
    if err := godotenv.Load(); err != nil {
        log.Println("⚠️  .env not found, using system env")
    }

    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("❌ DB_URL not set")
    }

    var err error
    DB, err = sqlx.Connect("postgres", dbURL)
    if err != nil {
        log.Fatalf("❌ Failed to connect to DB: %v", err)
    }

    fmt.Println("✅ Connected to PostgreSQL")

    schema := `
    CREATE TABLE IF NOT EXISTS clusters (
        id TEXT PRIMARY KEY,
        name TEXT NOT NULL,
        url TEXT NOT NULL,
        token TEXT NOT NULL
    );`
    DB.MustExec(schema)
}
