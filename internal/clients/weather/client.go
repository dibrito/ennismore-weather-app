package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	config "github.com/dibrito/ennismore-weather-app/config"
	"github.com/dibrito/ennismore-weather-app/internal/controller"
	"github.com/dibrito/ennismore-weather-app/pkg/logging"
	"github.com/dibrito/ennismore-weather-app/pkg/model"
	"go.uber.org/zap"
)

type Client struct {
	URL     string
	Timeout int
	Client  *http.Client
}

func New(c config.WeatherAPIConfig) *Client {
	return &Client{
		URL: c.URL,
		// Create an HTTP client with a timeout
		Client: &http.Client{
			Timeout: time.Duration(time.Duration(c.Timeout) * time.Second),
		},
	}
}

func (c *Client) GetForecast(ctx context.Context, lat, long string) ([]model.Period, error) {
	var result []model.Period
	// fetch resp
	resp, err := c.FetchForecastURL(ctx, lat, long)
	if err != nil {
		return result, err
	}

	// fetch model.periods
	return c.FetchForecastPeriods(ctx, resp.Properties.ForecastURL)
}

func (c *Client) FetchForecastURL(ctx context.Context, lat, long string) (model.WeatherPointsResponse, error) {
	logger := logging.GetLoggerFromContext(ctx)
	var result model.WeatherPointsResponse
	// bind URI with query param
	URI := fmt.Sprintf("%s/points/%s,%s", c.URL, lat, long)

	logger.Info("points URL", zap.String("url", URI))
	req, err := http.NewRequest(http.MethodGet, URI, nil)
	if err != nil {
		return result, err
	}

	req = req.WithContext(ctx)
	resp, err := c.Client.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusFound {
		return result, controller.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("unexpected status code: %d", resp.StatusCode)

	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to read response body: %v", err)
	}

	return result, nil
}

func (c *Client) FetchForecastPeriods(ctx context.Context, pointsURL string) ([]model.Period, error) {
	logger := logging.GetLoggerFromContext(ctx)
	var result model.ForecastResponse
	URI := pointsURL
	logger.Info("forecast URL", zap.String("url", URI))

	req, err := http.NewRequest(http.MethodGet, URI, nil)
	if err != nil {
		return []model.Period{}, err
	}
	req = req.WithContext(ctx)
	resp, err := c.Client.Do(req)
	if err != nil {
		return []model.Period{}, fmt.Errorf("failed to fetch data: %v", err)
	}
	defer resp.Body.Close()

	// TODO: check if theres any chance the API returning not found here. I don't think so!
	if resp.StatusCode == http.StatusFound {
		return []model.Period{}, controller.ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return []model.Period{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return []model.Period{}, fmt.Errorf("failed to read response body: %v", err)
	}

	return result.Properties.Periods, nil
}
