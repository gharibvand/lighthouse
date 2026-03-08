package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	StoragePath string
	RedisURL    string
}

func Load() *Config {
	// Load .env file (ignore error if file doesn't exist)
	// Try multiple strategies to find .env file
	
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	
	envPaths := []string{
		".env",                                    // Relative to current dir
		filepath.Join(wd, ".env"),                // Absolute path from working dir
	}
	
	// Try to find .env in current and parent directories (up to 3 levels)
	dir := wd
	for i := 0; i < 3; i++ {
		envPath := filepath.Join(dir, ".env")
		if _, statErr := os.Stat(envPath); statErr == nil {
			// File exists, add to beginning of list (priority)
			envPaths = append([]string{envPath}, envPaths...)
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	
	// Remove duplicates while preserving order
	seen := make(map[string]bool)
	uniquePaths := []string{}
	for _, path := range envPaths {
		if !seen[path] {
			seen[path] = true
			uniquePaths = append(uniquePaths, path)
		}
	}
	
	var loaded bool
	var loadedPath string
	for _, path := range uniquePaths {
		// Check if file exists first
		if _, err := os.Stat(path); err != nil {
			continue
		}
		
		// Try to load
		if err := godotenv.Load(path); err == nil {
			loaded = true
			loadedPath = path
			break
		} else {
			log.Printf("Failed to load .env from %s: %v", path, err)
		}
	}
	
	if loaded {
		absPath, _ := filepath.Abs(loadedPath)
		log.Printf("✓ Loaded .env file from: %s", absPath)
	} else {
		log.Println("⚠ No .env file found, using environment variables or defaults")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://lighthouse_user:lighthouse_password@localhost:5432/lighthouse?sslmode=disable"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		StoragePath: getEnv("STORAGE_PATH", "./storage"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
