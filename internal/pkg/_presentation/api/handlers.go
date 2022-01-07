package api

import (
	"compress/flate"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/diwise/api-smartwater/internal/pkg/application"
	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

func RegisterHandlers(r chi.Router, app application.Application, log zerolog.Logger) error {
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	// Enable gzip compression for ngsi-ld responses
	compressor := middleware.NewCompressor(flate.DefaultCompression, "application/json", "application/ld+json")
	r.Use(compressor.Handler)
	r.Use(middleware.Logger)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/ngsi-ld/v1/entities", NewCreateWaterConsumptionObservedHandler(log, app))

	return nil
}

func NewCreateWaterConsumptionObservedHandler(log zerolog.Logger, app application.Application) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fiwareWCO := fiware.WaterConsumptionObserved{}

		jsonBytes, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(jsonBytes, &fiwareWCO)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal request body into fiware entity")
		}

		timestamp, err := time.Parse(time.RFC3339, fiwareWCO.WaterConsumption.ObservedAt)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse time from string")
		}

		err = app.UpdateWaterConsumption(fiwareWCO.ID, fiwareWCO.WaterConsumption.Value, timestamp)
		if err != nil {
			log.Error().Err(err).Msgf("failed to create new entry in database because: %s", err.Error())
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json+ld")
		w.Write([]byte("creating new water consumption, in theory"))
	})
}
