package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-chi/httplog"
	"github.com/rs/zerolog/log"
	"github.com/sundsvall/api-smartwater/internal/pkg/_presentation/api"
	"github.com/sundsvall/api-smartwater/internal/pkg/application"
	"github.com/sundsvall/api-smartwater/internal/pkg/infrastructure/repositories/database"
)

func main() {
	serviceName := "api-smartwater"

	logger := log.With().Str("service", strings.ToLower(serviceName)).Logger()
	logger.Info().Msg("starting up ...")

	db, err := database.NewDatabaseConnection(database.NewPostgreSQLConnector(logger))
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to database, shutting down... ")
	}

	app := application.NewApplication(db, logger, serviceName)

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(
		httplog.NewLogger(serviceName, httplog.Options{
			JSON: true,
		}),
	))
	api.RegisterHandlers(r, app, logger)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("starting to listen for connections")

	log.Log().Str("Starting api-opendata on port:%s", port)
	err = http.ListenAndServe(":"+port, r)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to listen for connections")
	}
}
