package main

import (
	"compress/flate"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)

func main() {
	serviceName := "api-smartwater"

	logger := log.With().Str("service", strings.ToLower(serviceName)).Logger()
	logger.Info().Msg("starting up ...")

	router := createRequestRouter(serviceName)

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}

	log.Info().Str("port", port).Msg("starting to listen for connections")

	err := http.ListenAndServe(":"+port, router.impl)
	log.Fatal().Err(err).Msg("failed to listen for connections")
}

//RequestRouter wraps the concrete router implementation
type RequestRouter struct {
	impl *chi.Mux
}

//Get accepts a pattern that should be routed to the handlerFn on a GET request
func (router *RequestRouter) Get(pattern string, handlerFn http.HandlerFunc) {
	router.impl.Get(pattern, handlerFn)
}

//Post accepts a pattern that should be routed to the handlerFn on a POST request
func (router *RequestRouter) Post(pattern string, handlerFn http.HandlerFunc) {
	router.impl.Post(pattern, handlerFn)
}

func createRequestRouter(serviceName string) *RequestRouter {
	router := &RequestRouter{impl: chi.NewRouter()}

	router.impl.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)

	// Enable gzip compression for ngsi-ld responses
	compressor := middleware.NewCompressor(flate.DefaultCompression, "application/json", "application/ld+json")
	router.impl.Use(compressor.Handler)
	router.impl.Use(middleware.Logger)
	router.impl.Use(httplog.RequestLogger(
		httplog.NewLogger(serviceName, httplog.Options{
			JSON: true,
		}),
	))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Post("/ngsi-ld/v1/entities", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return router
}
