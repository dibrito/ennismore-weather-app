package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	serviceConfig "github.com/dibrito/ennismore-weather-app/config"
	"github.com/dibrito/ennismore-weather-app/internal/clients/openstreetmap"
	"github.com/dibrito/ennismore-weather-app/internal/clients/weather"
	repository "github.com/dibrito/ennismore-weather-app/internal/repository"
	"github.com/dibrito/ennismore-weather-app/pkg/logging"
	"github.com/dibrito/ennismore-weather-app/pkg/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// TODO: proper create integration tests

func TestIntegration(t *testing.T) {
	// this was a initial tests for the clients
	openstreetmapClient := openstreetmap.New(serviceConfig.OpenstreetmapAPIConfig{
		URL:     "https://nominatim.openstreetmap.org/search",
		Timeout: 30,
	})
	require.NotNil(t, openstreetmapClient)

	logger := zaptest.NewLogger(t)
	// Create a context with the logger
	ctx := context.WithValue(context.Background(), logging.LoggetCtxKey{}, logger)
	res, err := openstreetmapClient.GetLocation(ctx, "new york")
	require.NoError(t, err)
	if len(res) == 0 {
		t.Errorf("want:>0 got:%v", len(res))
	}
	bs, _ := json.Marshal(res)

	weatherClient := weather.New(serviceConfig.WeatherAPIConfig{
		URL:     "https://api.weather.gov",
		Timeout: 30,
	})
	require.NotNil(t, weatherClient)

	resPeriods, err := weatherClient.GetForecast(ctx, res[0].Lat, res[0].Lon)
	require.NoError(t, err)

	bs, _ = json.Marshal(resPeriods)
	logger.Info(string(bs))

	// but on integration tests we should spin up the whole solution like:
	// start the HTTP server
	// create reqs and check responses
}

func TestCache(t *testing.T) {
	cache := repository.New()
	cache.PutPeriods("city", "start", model.Period{})
	fmt.Println(cache.GetCache())
}
