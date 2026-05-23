package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := loadConfig()

	db, err := sql.Open("pgx", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database open failed: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("database connection failed: %v", err)
	}

	a := newApp(cfg, db)
	defer a.producer.Close()

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           a.routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("user-service listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server failed: %v", err)
	}
}
