package main

import (
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/sundsvall/api-smartwater/internal/pkg/application"
	"github.com/sundsvall/api-smartwater/internal/pkg/infrastructure/repositories/database"
)

func main() {
	serviceName := "api-smartwater"

	logger := log.With().Str("service", strings.ToLower(serviceName)).Logger()
	logger.Info().Msg("starting up ...")

	r := chi.NewRouter()

	db, err := database.NewDatabaseConnection(database.NewSQLiteConnector(logger))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database, shutting down... ")
	}

	app := application.NewApplication(r, db, logger, serviceName)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("starting to listen for connections")

	err = app.Start(port)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen for connections")
	}
}
