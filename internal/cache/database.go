package cache

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"
)

//go:embed schema/*.sql
var schemaFS embed.FS

func getDBPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("could not locate homne directory: %v", err)
		os.Exit(1)
	}

	var dbFolder string

	if runtime.GOOS == "windows" {
		dbFolder = filepath.Join(homeDir, "AppData", "Local", "prayertimes")
	} else {
		dbFolder = filepath.Join(homeDir, ".cache", "prayertimes")
	}

	if err := os.MkdirAll(dbFolder, 0o755); err != nil {
		fmt.Printf("error creating cache directory: %v", err)
		os.Exit(1)
	}

	return filepath.Join(dbFolder, "prayertimes.sqlite")
}

func DBExists() bool {
	dbPath := getDBPath()
	if _, err := os.Stat(dbPath); err == nil {
		return true
	}
	return false
}

// EnsureDB ensures the database exists and is up to date
func EnsureDB() (*sql.DB, error) {
	dbPath := getDBPath()

	dbFolder := filepath.Dir(dbPath)

	if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(dbFolder, 0o755); err != nil {
			return nil, fmt.Errorf("could not create cache directory: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if !DBExists() {
		fmt.Println("Database does not exist. Initializing...")

		goose.SetBaseFS(schemaFS)
		if err := goose.SetDialect("sqlite3"); err != nil {
			fmt.Printf("error setting dialect: %v", err)
			os.Exit(1)
		}
		if err := goose.Up(db, "schema"); err != nil {
			fmt.Printf("error running migrations: %v", err)
			os.Exit(1)
		}
	} else {
		// goose.SetBaseFS(schemaFS)
		//
		// if err := goose.SetDialect("sqlite3"); err != nil {
		// 	fmt.Printf("error setting dialect: %v", err)
		// 	os.Exit(1)
		// }
		//
		// current, err := goose.GetDBVersion(db)
		// if err != nil {
		// 	fmt.Printf("error getting db version; %v", err)
		// 	os.Exit(1)
		// }
		//
		// fmt.Printf("current db version: %d\n", current)
	}

	return db, nil
}
