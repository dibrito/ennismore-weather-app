package controller

import (
	"context"
	"errors"
	"testing"
	"time"

	openStreetMapAPIMock "github.com/dibrito/ennismore-weather-app/gen/mock/clients/openstreetmap"
	weatherAPIMock "github.com/dibrito/ennismore-weather-app/gen/mock/clients/weather"
	repositoryMock "github.com/dibrito/ennismore-weather-app/gen/mock/repository/memory"
	"github.com/dibrito/ennismore-weather-app/pkg/logging"
	"github.com/dibrito/ennismore-weather-app/pkg/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
)

var (
	location = model.Location{
		Lat: "123",
		Lon: "-456",
	}
)

func TestGetForecast(t *testing.T) {
	originalNowFunc := nowFunc
	defer func() { nowFunc = originalNowFunc }()

	// Set nowFunc to return a specific time
	fakeTime := time.Date(2024, 9, 23, 8, 0, 0, 0, time.UTC)
	nowFunc = func() time.Time {
		return fakeTime
	}
	tcs := []struct {
		name          string
		cities        []string
		checkResponse func(t *testing.T, got model.WeatherForecast, err error)
		setupMocks    func(t *testing.T, mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
			weatherMock *weatherAPIMock.MockWeatherGateway, cacheMock *repositoryMock.MockRepository)
	}{
		{
			name:   "when location/periods in cache should return cached data",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{
					Forecast: []model.Forecast{
						{
							Name: "london",
							Detail: []model.Detail{
								{
									StartTime:   nowFunc(),
									EndTime:     nowFunc().Add(2 * time.Hour),
									Description: "gray",
								},
								{
									StartTime:   nowFunc().Add(24 * time.Hour),
									EndTime:     nowFunc().Add(26 * time.Hour),
									Description: "gray",
								},
								{
									StartTime:   nowFunc().Add(48 * time.Hour),
									EndTime:     nowFunc().Add(50 * time.Hour),
									Description: "gray",
								},
							},
						},
					},
				}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(0)
				weatherMock.EXPECT().GetForecast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				cacheMock.EXPECT().GetLocation("london").Return(
					location,
					true).Times(1)
				setGetPeriodCalls(cacheMock)
			},
		},
		{
			name:   "when location not in cache and client call fails should skip city",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				cacheMock.EXPECT().GetLocation("london").Return(
					model.Location{},
					false).Times(1)
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(1).Return(
					[]model.Location{}, errors.New("client-err"),
				)
			},
		},
		{
			name:   "when location not in cache and client call is OK but location array is empty should skip city",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				cacheMock.EXPECT().GetLocation("london").Return(
					model.Location{},
					false).Times(1)
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(1).Return(
					[]model.Location{}, nil,
				)
			},
		},
		{
			name:   "when location not in cache and client call is OK should add location to cache and return forecast",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{
					Forecast: []model.Forecast{
						{
							Name: "london",
							Detail: []model.Detail{
								{
									StartTime:   nowFunc(),
									EndTime:     nowFunc().Add(2 * time.Hour),
									Description: "gray",
								},
								{
									StartTime:   nowFunc().Add(24 * time.Hour),
									EndTime:     nowFunc().Add(26 * time.Hour),
									Description: "gray",
								},
								{
									StartTime:   nowFunc().Add(48 * time.Hour),
									EndTime:     nowFunc().Add(50 * time.Hour),
									Description: "gray",
								},
							},
						},
					},
				}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				cacheMock.EXPECT().GetLocation("london").Return(
					model.Location{},
					false).Times(1)
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(1).Return(
					[]model.Location{location}, nil,
				)
				// add location to cache
				cacheMock.EXPECT().PutLocation("london", location)

				setGetPeriodCalls(cacheMock)
			},
		},
		{
			name:   "when periods not in cache and call client is ok should add periods to cache return forecast",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{
					Forecast: []model.Forecast{
						{
							Name: "london",
							Detail: []model.Detail{
								{
									StartTime:   nowFunc(),
									EndTime:     nowFunc().Add(2 * time.Hour),
									Description: "gray",
								},
								{
									StartTime:   nowFunc().Add(24 * time.Hour),
									EndTime:     nowFunc().Add(26 * time.Hour),
									Description: "gray",
								},
								{
									StartTime:   nowFunc().Add(48 * time.Hour),
									EndTime:     nowFunc().Add(50 * time.Hour),
									Description: "gray",
								},
							},
						},
					},
				}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(0)
				cacheMock.EXPECT().GetLocation("london").Return(
					location,
					true).Times(1)
				cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
					model.Period{}, false).Times(3)
				weatherMock.EXPECT().GetForecast(gomock.Any(), location.Lat, location.Lon).Times(1).Return(
					[]model.Period{
						{
							StartTime:   nowFunc(),
							EndTime:     nowFunc().Add(2 * time.Hour),
							Description: "gray",
						},
						{
							StartTime:   nowFunc().Add(24 * time.Hour),
							EndTime:     nowFunc().Add(26 * time.Hour),
							Description: "gray",
						},
						{
							StartTime:   nowFunc().Add(48 * time.Hour),
							EndTime:     nowFunc().Add(50 * time.Hour),
							Description: "gray",
						},
					}, nil)
				p1 := model.Period{
					StartTime:   nowFunc(),
					EndTime:     nowFunc().Add(2 * time.Hour),
					Description: "gray",
				}
				p2 := model.Period{
					StartTime:   nowFunc().Add(24 * time.Hour),
					EndTime:     nowFunc().Add(26 * time.Hour),
					Description: "gray",
				}
				p3 := model.Period{
					StartTime:   nowFunc().Add(48 * time.Hour),
					EndTime:     nowFunc().Add(50 * time.Hour),
					Description: "gray",
				}
				cacheMock.EXPECT().PutPeriods("london", p1.StartTime.Format("2006-01-02"), p1).Times(1)
				cacheMock.EXPECT().PutPeriods("london", p2.StartTime.Format("2006-01-02"), p2).Times(1)
				cacheMock.EXPECT().PutPeriods("london", p3.StartTime.Format("2006-01-02"), p3).Times(1)
			},
		},
		{
			name:   "when periods not in cache and call client NOT ok should skip city",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(0)
				cacheMock.EXPECT().GetLocation("london").Return(
					location,
					true).Times(1)
				cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
					model.Period{}, false).Times(3)
				weatherMock.EXPECT().GetForecast(gomock.Any(), location.Lat, location.Lon).Times(1).Return(
					[]model.Period{}, errors.New("client-error"))
			},
		},
		{
			name:   "when periods not matching days should skip city",
			cities: []string{"london"},
			checkResponse: func(t *testing.T, got model.WeatherForecast, err error) {
				require.NoError(t, err)
				want := model.WeatherForecast{}
				require.Equal(t, want, got)
			},
			setupMocks: func(
				t *testing.T,
				mapMock *openStreetMapAPIMock.MockOpenstreetmapperGateway,
				weatherMock *weatherAPIMock.MockWeatherGateway,
				cacheMock *repositoryMock.MockRepository) {
				mapMock.EXPECT().GetLocation(gomock.Any(), gomock.Any()).Times(0)
				weatherMock.EXPECT().GetForecast(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				cacheMock.EXPECT().GetLocation("london").Return(
					location,
					true).Times(1)
				cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
					model.Period{
						StartTime:   nowFunc().AddDate(0, 0, 10),
						EndTime:     nowFunc().AddDate(0, 0, 10).Add(time.Hour),
						Description: "gray",
					},
					true).Times(1)
				cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
					model.Period{
						StartTime:   nowFunc().AddDate(0, 0, 11),
						EndTime:     nowFunc().AddDate(0, 0, 11).Add(time.Hour),
						Description: "gray",
					},
					true).Times(1)
				cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
					model.Period{
						StartTime:   nowFunc().AddDate(0, 0, 12),
						EndTime:     nowFunc().AddDate(0, 0, 12).Add(time.Hour),
						Description: "gray",
					},
					true).Times(1)
			},
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repoMock := repositoryMock.NewMockRepository(ctrl)
			openStreetMapAPIMock := openStreetMapAPIMock.NewMockOpenstreetmapperGateway(ctrl)
			weatherAPIMock := weatherAPIMock.NewMockWeatherGateway(ctrl)

			weatherAppController := New(openStreetMapAPIMock, weatherAPIMock, repoMock)
			logger := zaptest.NewLogger(t)
			ctx := context.WithValue(context.Background(), logging.LoggetCtxKey{}, logger)

			tc.setupMocks(t, openStreetMapAPIMock, weatherAPIMock, repoMock)

			got, err := weatherAppController.GetForecast(ctx, tc.cities)
			tc.checkResponse(t, got, err)
		})
	}
}

func setGetPeriodCalls(cacheMock *repositoryMock.MockRepository) {
	cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
		model.Period{
			StartTime:   nowFunc(),
			EndTime:     nowFunc().Add(2 * time.Hour),
			Description: "gray",
		},
		true).Times(1)
	cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
		model.Period{
			StartTime:   nowFunc().Add(24 * time.Hour),
			EndTime:     nowFunc().Add(26 * time.Hour),
			Description: "gray",
		},
		true).Times(1)
	cacheMock.EXPECT().GetPeriods("london", gomock.Any()).Return(
		model.Period{
			StartTime:   nowFunc().Add(48 * time.Hour),
			EndTime:     nowFunc().Add(50 * time.Hour),
			Description: "gray",
		},
		true).Times(1)
}
