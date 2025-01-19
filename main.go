package main

import (
	"NablaFunctions/handlers"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	log.Info().Msg("Hello, Serverless User!")

	http.HandleFunc("/api/load", handlers.LoggingMiddleWare(handlers.LoadHandler))
	http.HandleFunc("/api/execute", handlers.LoggingMiddleWare(handlers.ExecuteHandler))

	log.Info().Msg("Server listening on port 8080...")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
