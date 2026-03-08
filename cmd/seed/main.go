package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"lighthouse/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	// Connect to database
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	fmt.Println("✓ Connected to database")

	// Find seeder directory
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Try to find seeders directory (could be in root or cmd/seed)
	seederPaths := []string{
		filepath.Join(wd, "seeders"),
		filepath.Join(wd, "..", "seeders"),
		filepath.Join(wd, "..", "..", "seeders"),
	}

	var seederDir string
	for _, path := range seederPaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			seederDir = path
			break
		}
	}

	if seederDir == "" {
		log.Fatalf("Seeders directory not found. Tried: %v", seederPaths)
	}

	fmt.Printf("✓ Found seeders directory: %s\n", seederDir)

	// Read and execute all SQL files in seeders directory
	files, err := os.ReadDir(seederDir)
	if err != nil {
		log.Fatalf("Failed to read seeders directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}

		filePath := filepath.Join(seederDir, file.Name())
		fmt.Printf("\n📄 Executing: %s\n", file.Name())

		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("⚠ Failed to read %s: %v", file.Name(), err)
			continue
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			log.Printf("⚠ Failed to execute %s: %v", file.Name(), err)
			continue
		}

		fmt.Printf("✓ Successfully executed: %s\n", file.Name())
	}

	fmt.Println("\n✅ All seeders completed!")
}
