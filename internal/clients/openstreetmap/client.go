package openstreetmap

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

func New(c config.OpenstreetmapAPIConfig) *Client {
	return &Client{
		URL: c.URL,
		// Create an HTTP client with a timeout
		Client: &http.Client{
			Timeout: time.Duration(time.Duration(c.Timeout) * time.Second),
		},
	}
}

func (c *Client) GetLocation(ctx context.Context, parmas string) ([]model.Location, error) {
	logger := logging.GetLoggerFromContext(ctx)

	var result []model.Location
	// bind URI with query param
	URI := c.URL

	// create request
	req, err := http.NewRequest(http.MethodGet, URI, nil)
	if err != nil {
		return result, err
	}

	// bind req with ctx
	req = req.WithContext(ctx)
	values := req.URL.Query()
	values.Add("q", parmas)
	values.Add("format", "json")
	req.URL.RawQuery = values.Encode()

	logger.Info("calling URL", zap.String("url", req.URL.String()))

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
