package respository

import (
	"sync"

	"github.com/dibrito/ennismore-weather-app/pkg/model"
)

type Repository interface {
	GetLocation(city string) (model.Location, bool)
	PutLocation(city string, location model.Location)
	GetPeriods(city, startTime string) (model.Period, bool)
	PutPeriods(city, startTime string, periods model.Period)
	GetCache() model.CacheResponse
}

// Repository defines a in-memory weather-app repository.
type repository struct {
	sync.RWMutex
	Location map[string]model.Location
	Periods  map[string]map[string]model.Period
}

// New creates a new memory repository.
func New() *repository {
	return &repository{
		Location: make(map[string]model.Location, 0),
		Periods:  make(map[string]map[string]model.Period, 0),
	}
}

// GetLocation retrieves the Location by city name.
func (r *repository) GetLocation(city string) (model.Location, bool) {
	r.RLock()
	defer r.RUnlock()

	location, ok := r.Location[city]
	return location, ok
}

// PutLocation stores a Location by city name.
func (r *repository) PutLocation(city string, location model.Location) {
	r.Lock()
	defer r.Unlock()

	r.Location[city] = location
}

// GetPeriods retrieves the forecast periods by city name.
func (r *repository) GetPeriods(city, startTime string) (model.Period, bool) {
	r.RLock()
	defer r.RUnlock()

	periods, ok := r.Periods[city]
	if !ok {
		return model.Period{}, ok
	}
	p, ok := periods[startTime]
	return p, ok
}

// PutPeriods stores forecast periods for a city.
func (r *repository) PutPeriods(city, startTime string, periods model.Period) {
	r.Lock()
	defer r.Unlock()
	if _, ok := r.Periods[city]; !ok {
		r.Periods[city] = make(map[string]model.Period)
	}
	r.Periods[city][startTime] = periods
}

// GetLocation retrieves the Location by city name.
func (r *repository) GetCache() model.CacheResponse {
	r.RLock()
	defer r.RUnlock()
	return model.CacheResponse{
		Location: r.Location,
		Periods:  r.Periods,
	}
}
