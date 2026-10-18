package main

import (
	_ "api/docs"
	"api/internal/config"
	"api/internal/health"
	"api/internal/postgres"
	"api/internal/spending"
	"api/internal/user"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

// @title Telegram Cash Tracker API
// @version 1.0
// @description API for tracking family expenses.
// @BasePath /
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
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("/docs/doc.json")))

	// users
	userHandler := user.NewHandler(user.NewRepository(db))
	r.Post("/users", userHandler.Create)

	spendingRepository := spending.NewRepository(db)
	spendingHandler := spending.NewHandler(spendingRepository)
	r.Post("/spendings", spendingHandler.Create)
	r.Get("/spendings", spendingHandler.List)
	r.Patch("/spendings/{id}", spendingHandler.Update)
	r.Delete("/spendings/{id}", spendingHandler.Delete)

	adminSpendingHandler := spending.NewAdminHandler(spendingRepository)
	r.Route("/admin", func(r chi.Router) {
		r.Post("/spendings", adminSpendingHandler.Create)
		r.Get("/spendings", adminSpendingHandler.List)
		r.Get("/spendings/{id}", adminSpendingHandler.Get)
		r.Patch("/spendings/{id}", adminSpendingHandler.Update)
		r.Delete("/spendings/{id}", adminSpendingHandler.Delete)
	})

	fmt.Printf("server is working on %s\n", config.Port)
	http.ListenAndServe(":"+config.Port, r)
}
