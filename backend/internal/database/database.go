// Package database provee una conexión singleton a PostgreSQL via GORM.
package database

import (
	"fmt"
	"sync"

	"github.com/Alexisjar91/POS/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// Get retorna la instancia singleton de *gorm.DB.
// Conecta a PostgreSQL la primera vez que se llama usando config.Get().
// Panic si no puede conectar — la app no puede funcionar sin DB.
func Get() *gorm.DB {
	once.Do(func() {
		cfg := config.Get()
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
		)
		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}
	})
	return db
}
