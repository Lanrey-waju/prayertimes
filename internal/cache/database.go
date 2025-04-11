package cache

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
	_ "modernc.org/sqlite"
)

type OSProvider interface {
	GetOS() string
	GetUserHomeDir() (string, error)
	MkdirAll(path string, perm os.FileMode) error
}

type DefaultOSProvider struct{}

func NewDefaultOSProvider() OSProvider {
	return DefaultOSProvider{}
}

func (p DefaultOSProvider) GetOS() string {
	return runtime.GOOS
}

func (p DefaultOSProvider) GetUserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (p DefaultOSProvider) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

//go:embed schema/*.sql
var schemaFS embed.FS

func GetDBPath(osProvider OSProvider) (string, error) {
	homeDir, err := osProvider.GetUserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get user home directory: %v", err)
	}

	var dbFolder string

	if osProvider.GetOS() == "windows" {
		dbFolder = filepath.Join(homeDir, "AppData", "Local", "prayertimes")
	} else {
		dbFolder = filepath.Join(homeDir, ".cache", "prayertimes")
	}

	if err := osProvider.MkdirAll(dbFolder, 0o755); err != nil {
		return "", fmt.Errorf("could not create cache directory: %v", err)
	}

	return filepath.Join(dbFolder, "prayertimes.sqlite"), nil
}

func DBExists(osProvider OSProvider) bool {
	dbPath, err := GetDBPath(osProvider)
	if err != nil {
		fmt.Printf("error getting db path: %v", err)
	}
	if _, err := os.Stat(dbPath); err == nil {
		return true
	}
	return false
}

// EnsureDB ensures the database exists and is up to date
func EnsureDB(osProvider OSProvider) (*sql.DB, error) {
	dbPath, _ := GetDBPath(osProvider)

	dbFolder := filepath.Dir(dbPath)

	if _, err := os.Stat(dbFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(dbFolder, 0o755); err != nil {
			return nil, fmt.Errorf("could not create cache directory: %v", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	if !DBExists(osProvider) {
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
