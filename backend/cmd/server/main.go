package main

import (
	"log"
	"os"

	"github.com/flynnzhang/planning/backend/internal/database"
	"github.com/flynnzhang/planning/backend/internal/handler"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "rdc.db"
	}

	db, err := database.New(dsn)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	addr := os.Getenv("LISTEN_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	r := handler.SetupRouter(db)
	log.Printf("starting server on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
