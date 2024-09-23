package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	controllerMock "github.com/dibrito/ennismore-weather-app/gen/mock/controller"
	"github.com/dibrito/ennismore-weather-app/internal/controller"
	"github.com/dibrito/ennismore-weather-app/pkg/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var want = model.WeatherForecast{
	Forecast: []model.Forecast{
		{
			Name: "london",
			Detail: []model.Detail{
				{
					StartTime:   time.Now().UTC(),
					EndTime:     time.Now().UTC().Add(2 * time.Hour),
					Description: "gray",
				},
				{
					StartTime:   time.Now().UTC().Add(24 * time.Hour),
					EndTime:     time.Now().UTC().Add(26 * time.Hour),
					Description: "gray",
				},
				{
					StartTime:   time.Now().UTC().Add(48 * time.Hour),
					EndTime:     time.Now().UTC().Add(50 * time.Hour),
					Description: "gray",
				},
			},
		},
	},
}

func TestGetForecast(t *testing.T) {
	tcs := []struct {
		name          string
		queryParams   string
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
		setupMock     func(mock *controllerMock.MockServiceController)
	}{
		{
			name:        "when empty query params should return BAD REQUEST",
			queryParams: "",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Result().StatusCode, http.StatusBadRequest)
			},
			setupMock: func(mock *controllerMock.MockServiceController) {
			},
		},
		{
			name:        "when error should return Status Internal Server Error",
			queryParams: "?city=london",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Result().StatusCode, http.StatusInternalServerError)
			},
			setupMock: func(mock *controllerMock.MockServiceController) {
				mock.EXPECT().GetForecast(gomock.Any(), []string{"london"}).Return(want, errors.New("service-error")).Times(1)
			},
		},
		{
			name:        "when error not found should return Status Not Found",
			queryParams: "?city=london",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, recorder.Result().StatusCode, http.StatusNotFound)
			},
			setupMock: func(mock *controllerMock.MockServiceController) {
				mock.EXPECT().GetForecast(gomock.Any(), []string{"london"}).Return(want, controller.ErrNotFound).Times(1)
			},
		},
		{
			name:        "when no error should return forecast",
			queryParams: "?city=london",
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var got model.WeatherForecast
				err := json.NewDecoder(recorder.Body).Decode(&got)
				require.NoError(t, err)
				// want :=
				require.Equal(t, want, got)
			},
			setupMock: func(mock *controllerMock.MockServiceController) {
				mock.EXPECT().GetForecast(gomock.Any(), []string{"london"}).Return(want, nil).Times(1)
			},
		},
	}
	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			serviceControllerMock := controllerMock.NewMockServiceController(ctrl)
			tc.setupMock(serviceControllerMock)
			handler := New(serviceControllerMock)

			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/weather%v", tc.queryParams)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			// http.DefaultServeMux.ServeHTTP(recorder, req)
			handler.GetForecast(recorder, req)
			tc.checkResponse(t, recorder)

			// handler.GetForecast(recorder, req)

		})
	}
}
