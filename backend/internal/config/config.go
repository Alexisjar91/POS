// Package config provee una configuración singleton cargada desde .env.
package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// Config contiene todas las variables de entorno del sistema.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

var (
	instance *Config
	once     sync.Once
)

// Get retorna la instancia singleton de Config.
// Carga .env la primera vez que se llama. No falla si .env no existe
// (las variables pueden venir del entorno del sistema).
func Get() *Config {
	once.Do(func() {
		// Intenta cargar .env, ignora error si no existe
		_ = godotenv.Load()

		instance = &Config{
			DBHost:     getEnv("DB_HOST", "localhost"),
			DBPort:     getEnv("DB_PORT", "5432"),
			DBUser:     getEnv("DB_USER", "pos_user"),
			DBPassword: getEnv("DB_PASSWORD", "pos_password"),
			DBName:     getEnv("DB_NAME", "pos_db"),
			ServerPort: getEnv("SERVER_PORT", "8080"),
		}
	})
	return instance
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
