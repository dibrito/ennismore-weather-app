package model

import "time"

//	ennismore-weather-app response:
//
// WeatherForecast represent the main forecast for weather-app API.
type WeatherForecast struct {
	Forecast []Forecast `json:"forecast"`
}

// Forecast represent the forecast for a city
type Forecast struct {
	Name   string   `json:"name"`
	Detail []Detail `json:"detail"`
}

// Detail represent the inner details of a forecast
type Detail struct {
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"description"`
}

//	gateways responses:
//
// Location represent the response for OpenStreetMap API requests.
type Location struct {
	PlaceID     int    `json:"place_id"`
	Licence     string `json:"licence"`
	OsmType     string `json:"osm_type"`
	OsmID       int    `json:"osm_id"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
	Class       string `json:"class"`
	Type        string `json:"type"`
}

// WeatherPointsResponse represent the response for Weather API requests.
type WeatherPointsResponse struct {
	Properties WeatherProperties `json:"properties"`
}

// WeatherProperties represent the forecast URL for a pair of points.
type WeatherProperties struct {
	ForecastURL string `json:"forecast"`
}

// ForecastResponse represent the response for calling forcast URL for a pair of points.
type ForecastResponse struct {
	Properties ForecastProperties `json:"properties"`
}

// ForecastProperties represent the details of a forecast for a pair of points.
type ForecastProperties struct {
	Periods []Period `json:"periods"`
}

// Period represents the forecast detail.
type Period struct {
	StartTime   time.Time `json:"startTime"`
	EndTime     time.Time `json:"endTime"`
	Description string    `json:"detailedForecast"`
}

type CacheResponse struct {
	Location map[string]Location          `json:"locations"`
	Periods  map[string]map[string]Period `json:"periods"`
}
