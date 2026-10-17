package main

import (
	"api/internal/health"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func main() {

	envFile := os.Getenv("ENV_FILE")

	if envFile != "" {
		if err := godotenv.Load(envFile); err != nil {
			log.Fatalf("failed load env %v", err)
		}
	}

	port := os.Getenv("API_PORT")

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", health.Ping)

	fmt.Printf("Server is working on %s\n", port)
	http.ListenAndServe(":"+port, r)
}
