package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/dibrito/ennismore-weather-app/internal/controller"
	"github.com/dibrito/ennismore-weather-app/pkg/logging"
	"github.com/dibrito/ennismore-weather-app/pkg/model"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"
)

var contentTypeKey = "Content-type"
var contentTypeValue = "application/json"

// Handler defines weather-app HTTP handler.
type Handler struct {
	ctrl ServiceController
}

type ServiceController interface {
	GetCache() model.CacheResponse
	GetForecast(ctx context.Context, cities []string) (model.WeatherForecast, error)
}

// New creates a new movie metadata HTTP handler.
func New(ctrl ServiceController) *Handler {
	return &Handler{ctrl}
}

// GetForecast handles GET /weather requests.
func (h *Handler) GetForecast(w http.ResponseWriter, req *http.Request) {
	logger := logging.GetLoggerFromContext(req.Context())
	logger = logger.With(zap.String("URI", req.RequestURI))

	citiesQuery := req.URL.Query().Get("city")
	if citiesQuery == "" {
		http.Error(w, "Please provide a list of cities", http.StatusBadRequest)
		return
	}
	logger = logger.With(zap.String("cities_query", citiesQuery))

	ctx := req.Context()
	params, err := parseCities(citiesQuery)
	if err != nil {
		logger.Error("unable to parse cites", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	m, err := h.ctrl.GetForecast(ctx, params)
	if err != nil && errors.Is(err, controller.ErrNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		logger.Error("unable retrieve forecast", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		logger.Error("unable to parse response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// Function to break down the city query string into a slice of strings
func parseCities(query string) ([]string, error) {
	cities := strings.Split(query, ",")
	var decodedCities []string

	// Loop through cities, decode percent-encoding, and append to the slice
	for _, city := range cities {
		decodedCity, err := url.QueryUnescape(strings.TrimSpace(city)) // Decoding %20 and trimming spaces
		if err != nil {
			return nil, err
		}
		decodedCities = append(decodedCities, decodedCity)
	}

	return decodedCities, nil
}

// TODO: I'm not happy about the handler calling routes!
func (h *Handler) Routes(logger *zap.Logger) http.Handler {
	mux := chi.NewRouter()

	// specify who is allwoed to connect
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", contentTypeKey, "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// define HTTP routes
	mux.Use(middleware.Heartbeat("/health"))
	mux.With(LoggerInterceptor(logger)).Get("/weather", http.HandlerFunc(h.GetForecast))
	mux.With(LoggerInterceptor(logger)).Get("/cache", http.HandlerFunc(h.GetCache))

	return mux
}

// LoggerInterceptor adds tracing info to the logger and puts the logger into request Context
func LoggerInterceptor(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(contentTypeKey, contentTypeValue)
			ctx := context.WithValue(r.Context(), logging.LoggetCtxKey{}, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// TODO: please disegard this endpoint, this was used to test purpose only to act like get all from cache!
func (h *Handler) GetCache(w http.ResponseWriter, req *http.Request) {
	logger := logging.GetLoggerFromContext(req.Context())
	cache := h.ctrl.GetCache()
	w.Header().Set(contentTypeKey, contentTypeValue)
	if err := json.NewEncoder(w).Encode(cache); err != nil {
		logger.Error("unable to parse response", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
