package main

import (
	"api/internal/config"
	"api/internal/health"
	"api/internal/postgres"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUsername,
		config.PostgresPassword,
		config.PostgresDB,
	)

	db, err := postgres.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("postgres connected")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", health.Ping)

	fmt.Printf("server is working on %s\n", config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
