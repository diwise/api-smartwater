package application

import (
	"compress/flate"
	"encoding/json"
	"io"
	"net/http"

	"github.com/diwise/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/sundsvall/api-smartwater/internal/pkg/infrastructure/repositories/database"
)

type Application interface {
	Start(port string) error
}

type waterConsumptionApp struct {
	router chi.Router
	db     database.Datastore
	log    zerolog.Logger
}

func (w *waterConsumptionApp) Start(port string) error {
	w.log.Log().Str("Starting api-opendata on port:%s", port)
	return http.ListenAndServe(":"+port, w.router)
}

func NewApplication(r chi.Router, db database.Datastore, log zerolog.Logger, serviceName string) Application {
	return newWaterConsumptionApp(r, db, log, serviceName)
}

func newWaterConsumptionApp(r chi.Router, db database.Datastore, log zerolog.Logger, serviceName string) *waterConsumptionApp {
	w := &waterConsumptionApp{
		router: r,
		db:     db,
		log:    log,
	}

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	// Enable gzip compression for ngsi-ld responses
	compressor := middleware.NewCompressor(flate.DefaultCompression, "application/json", "application/ld+json")
	r.Use(compressor.Handler)
	r.Use(middleware.Logger)
	r.Use(httplog.RequestLogger(
		httplog.NewLogger(serviceName, httplog.Options{
			JSON: true,
		}),
	))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Post("/ngsi-ld/v1/entities", w.createNewWaterConsumption())

	return w
}

func (wca *waterConsumptionApp) createNewWaterConsumption() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fiwareWCO := fiware.WaterConsumptionObserved{}

		jsonBytes, _ := io.ReadAll(r.Body)
		err := json.Unmarshal(jsonBytes, &fiwareWCO)
		if err != nil {
			wca.log.Error().Err(err).Msg("failed to unmarshal request body into fiware entity")
		}

		_, err = wca.db.CreateWaterConsumption(&fiwareWCO)
		if err != nil {
			wca.log.Error().Err(err).Msgf("failed to create new entry in database because: %s", err.Error())
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json+ld")
		w.Write([]byte("creating new water consumption, in theory"))
	})
}
