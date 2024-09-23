package controller

import (
	"context"
	"errors"
	"time"

	respository "github.com/dibrito/ennismore-weather-app/internal/repository"
	"github.com/dibrito/ennismore-weather-app/pkg/logging"
	"github.com/dibrito/ennismore-weather-app/pkg/model"
	"go.uber.org/zap"
)

// ErrNotFound is returned when a requested record is not
// found.
// TODO: move to pkg
var ErrNotFound = errors.New("not found")

// nowFunc will get the now time when GetForecast is executed.
var nowFunc = time.Now

type OpenstreetmapperGateway interface {
	GetLocation(ctx context.Context, city string) ([]model.Location, error)
}

type WeatherGateway interface {
	GetForecast(ctx context.Context, lat, long string) ([]model.Period, error)
}

// Controller defines a metadata service controller.
type Controller struct {
	openStreetMapperClient OpenstreetmapperGateway
	weatherClient          WeatherGateway
	cacheRepository        respository.Repository
}

// New creates a weather-app service controller.
func New(openStreetMapperClient OpenstreetmapperGateway,
	weatherClient WeatherGateway, cache respository.Repository) *Controller {
	return &Controller{
		openStreetMapperClient: openStreetMapperClient,
		weatherClient:          weatherClient,
		cacheRepository:        cache,
	}
}

// Get returns movie metadata by id.
func (c *Controller) GetForecast(ctx context.Context, cities []string) (model.WeatherForecast, error) {
	logger := logging.GetLoggerFromContext(ctx)
	var result model.WeatherForecast

	// TODO probably don't need this
	// since the firs element of the API is "today" always
	now := nowFunc().UTC()
	days := []time.Time{now, now.AddDate(0, 0, 1), now.AddDate(0, 0, 2)}

	// get points for each city
	for _, city := range cities {
		var location model.Location
		locationCache, ok := c.cacheRepository.GetLocation(city)
		location = locationCache
		if !ok {
			logger.Info("location not found in cache, calling client",
				zap.String("location", city))
			locationClient, err := c.openStreetMapperClient.GetLocation(ctx, city)
			if err != nil {
				// here we won't fail the whole operation but log the failures
				// better approach can be discussed, e.g. we could return the response
				// also with failures, or provide some report of the operation.
				logger.Warn("unable to retrieve location",
					zap.String("location", city),
					zap.Error(err))
				continue
			}
			// location = locationClient
			if len(locationClient) == 0 {
				logger.Info("location not found",
					zap.String("location", city))
				continue
			}
			// add to cache
			location = locationClient[0]
			c.cacheRepository.PutLocation(city, location)
		}

		// for each city I need 3 periods: now + 2days
		// either I get from cache
		periods := c.getPeriodsFromCache(city, days)
		logger.Info("cach periods",
			zap.Any("periods", periods))
		if len(periods) != 3 {
			// or I query 3rd API
			logger.Info("periods not found in cache, calling client",
				zap.Any("days", days))
			periodsClient, err := c.weatherClient.GetForecast(ctx, location.Lat, location.Lon)
			if err != nil {
				logger.Warn("unable to retrieve forecast",
					zap.String("location", city),
					zap.String("lat", location.Lat),
					zap.String("log", location.Lon),
					zap.Error(err))
				continue
			}
			for _, p := range periodsClient {
				c.cacheRepository.PutPeriods(city, p.StartTime.Format("2006-01-02"), p)
			}
			periods = periodsClient
		}

		logger.Info("finnding forecasts details",
			zap.Any("days", days),
			zap.Any("periods", periods),
		)

		// here we need today's forecast and next 2 days
		details := findForecast(periods, days)
		if len(details) == 0 {
			logger.Warn("unable find forecast for time period(now+2days)",
				zap.String("location", city),
				zap.String("lat", location.Lat),
				zap.String("log", location.Lon))
			continue
		}

		result.Forecast = append(result.Forecast, model.Forecast{
			Name:   city,
			Detail: details,
		})
	}

	// fill cache

	return result, nil
}

// getPeriodsFromCache will get all posible periods: now:2days
func (c *Controller) getPeriodsFromCache(city string, days []time.Time) []model.Period {
	var result []model.Period
	for _, day := range days {
		if p, ok := c.cacheRepository.GetPeriods(city, day.Format("2006-01-02")); ok {
			result = append(result, p)
		}
	}
	return result
}

// findForecast finds forecast for today + 2 days.
func findForecast(periods []model.Period, days []time.Time) []model.Detail {
	var result []model.Detail
	for _, day := range days {
		for _, p := range periods {
			if matchDate(p.StartTime, day) {
				result = append(result, model.Detail(p))
				break
			}
		}
	}

	return result
}

func matchDate(a, b time.Time) bool {
	yearA, monthA, dayA := a.Date()
	yearB, monthB, dayB := b.Date()
	return yearA == yearB && monthA == monthB && dayA == dayB
}

// NOTE:this is test purpose only to get cache content during app execution, disregard this func.
func (c *Controller) GetCache() model.CacheResponse {
	return c.cacheRepository.GetCache()
}
