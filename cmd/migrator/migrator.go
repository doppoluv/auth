package main

import (
	"database/sql"
	"flag"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"

	"auth/internal/lib/logger"
)

func main() {
	log := logger.NewLogger()

	var storagePath, migrationPath, migrationTable string

	flag.StringVar(&storagePath, "storagepath", "", "path to the storage file")
	flag.StringVar(&migrationPath, "migrationpath", "", "path to the migrations directory")
	flag.StringVar(&migrationTable, "migrationtable", "migrations", "name of the migrations table")
	flag.Parse()

	if storagePath == "" || migrationPath == "" {
		flag.Usage()
		log.Fatalf("storage path and migrations path are required")
	}

	db, err := sql.Open("sqlite", storagePath)
	if err != nil {
		log.Fatalf("open sqlite database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ping sqlite database: %v", err)
	}

	goose.SetDialect("sqlite3")
	goose.SetTableName(migrationTable)

	if err := goose.Up(db, migrationPath); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	log.Printf("Migrations applied successfully")
}
