package api

import (
	"compress/flate"
	"net/http"

	context "github.com/diwise/api-smartwater/internal/pkg/_presentation/api/ngsi-ld/context"
	"github.com/diwise/api-smartwater/internal/pkg/application"
	"github.com/diwise/ngsi-ld-golang/pkg/ngsi-ld"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

func createContextRegistry(app application.Application, log zerolog.Logger) ngsi.ContextRegistry {
	contextRegistry := ngsi.NewContextRegistry()
	ctxSource := context.CreateSource(app, log)
	contextRegistry.Register(ctxSource)
	return contextRegistry
}

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

	ctxReg := createContextRegistry(app, log)

	r.Post("/ngsi-ld/v1/entities", ngsi.NewCreateEntityHandler(ctxReg))
	r.Get("/ngsi-ld/v1/entities", ngsi.NewQueryEntitiesHandler(ctxReg))

	return nil
}
