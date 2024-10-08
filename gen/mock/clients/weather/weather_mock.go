// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/dibrito/ennismore-weather-app/internal/controller (interfaces: WeatherGateway)
//
// Generated by this command:
//
//	mockgen -package weather_mock --destination=./gen/mock/clients/weather/weather_mock.go github.com/dibrito/ennismore-weather-app/internal/controller WeatherGateway
//

// Package weather_mock is a generated GoMock package.
package weather_mock

import (
	context "context"
	reflect "reflect"

	model "github.com/dibrito/ennismore-weather-app/pkg/model"
	gomock "go.uber.org/mock/gomock"
)

// MockWeatherGateway is a mock of WeatherGateway interface.
type MockWeatherGateway struct {
	ctrl     *gomock.Controller
	recorder *MockWeatherGatewayMockRecorder
}

// MockWeatherGatewayMockRecorder is the mock recorder for MockWeatherGateway.
type MockWeatherGatewayMockRecorder struct {
	mock *MockWeatherGateway
}

// NewMockWeatherGateway creates a new mock instance.
func NewMockWeatherGateway(ctrl *gomock.Controller) *MockWeatherGateway {
	mock := &MockWeatherGateway{ctrl: ctrl}
	mock.recorder = &MockWeatherGatewayMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWeatherGateway) EXPECT() *MockWeatherGatewayMockRecorder {
	return m.recorder
}

// GetForecast mocks base method.
func (m *MockWeatherGateway) GetForecast(arg0 context.Context, arg1, arg2 string) ([]model.Period, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetForecast", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.Period)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetForecast indicates an expected call of GetForecast.
func (mr *MockWeatherGatewayMockRecorder) GetForecast(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetForecast", reflect.TypeOf((*MockWeatherGateway)(nil).GetForecast), arg0, arg1, arg2)
}
