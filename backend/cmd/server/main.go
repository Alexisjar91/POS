package main

import (
	"log"
	"net/http"

	"github.com/Alexisjar91/POS/internal/config"
	"github.com/Alexisjar91/POS/internal/database"
	"github.com/Alexisjar91/POS/internal/users/infrastructure/persistence/postgres"
)

func main() {
	cfg := config.Get()
	db := database.Get()

	if err := postgres.RunMigrations(db); err != nil {
		log.Fatalf("migrations: %v", err)
	}
	if err := postgres.RunSeed(db); err != nil {
		log.Fatalf("seed: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Printf("POS API ready on :%s", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}
