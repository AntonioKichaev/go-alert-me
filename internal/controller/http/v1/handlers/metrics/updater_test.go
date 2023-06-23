package metrics_test

import (
	"github.com/antoniokichaev/go-alert-me/internal/usecase/repo/mocks"
	"github.com/antoniokichaev/go-alert-me/pkg/metrics"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	mockStore := mocks.NewKeeper(t)

	srv := getServer(mockStore)
	defer srv.Close()

	const _addCounter = "AddCounter"
	const _setGauge = "SetGauge"
	type mockStoreRequest struct {
		methodName string
		args       []any
		returnArgs []any
	}
	tt := map[string]struct {
		method      string
		targetURL   string
		statusCode  int
		contentType string
		mockStore   mockStoreRequest
		wantErr     error
	}{
		"add counter ": {
			method:      http.MethodPost,
			targetURL:   "/update/counter/1/2",
			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore: mockStoreRequest{
				methodName: _addCounter,
				args:       []any{&metrics.Counter{Name: "1", Value: 2}},
				returnArgs: []any{&metrics.Counter{Name: "1", Value: 2}, nil},
			},
		},
		"zero_value ": {
			method:      http.MethodPost,
			targetURL:   "/update/counter/1/",
			statusCode:  http.StatusNotFound,
			contentType: _contentTypeText,
		},
		"zero_metrics ": {
			method:      http.MethodPost,
			targetURL:   "/update/counter//5",
			statusCode:  http.StatusNotFound,
			contentType: _contentTypeText,
			mockStore:   mockStoreRequest{},
		},
		"unknown_metric_type ": {
			method:      http.MethodPost,
			targetURL:   "/update/xep/er/5",
			statusCode:  http.StatusBadRequest,
			contentType: "",
			mockStore:   mockStoreRequest{},
		},
		"unk ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusBadRequest,
			contentType: "",
			mockStore:   mockStoreRequest{methodName: _addCounter},
		},
		"negative_value ": {
			method:      http.MethodPost,
			targetURL:   "/update/counter/ram/-5",
			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore: mockStoreRequest{methodName: _addCounter,
				args:       []any{&metrics.Counter{Name: "ram", Value: int64(-5)}},
				returnArgs: []any{&metrics.Counter{Name: "ram", Value: int64(-5)}, nil},
			},
		},
		"negative_float_value ": {
			method:      http.MethodPost,
			targetURL:   "/update/counter/ram/-5.5",
			statusCode:  http.StatusBadRequest,
			contentType: "",
			mockStore:   mockStoreRequest{},
		},
		"simple_set_gauge ": {
			method:      http.MethodPost,
			targetURL:   "/update/gauge/ram/999.5999",
			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore: mockStoreRequest{
				methodName: _setGauge,
				args:       []any{&metrics.Gauge{Name: "ram", Value: 999.5999}},
				returnArgs: []any{&metrics.Gauge{Name: "ram", Value: 999.5999}, nil},
			},
		},
		"none_value_set_gauge ": {
			method:      http.MethodPost,
			targetURL:   "/update/gauge/ram/none",
			statusCode:  http.StatusBadRequest,
			contentType: "",
		},
		"incorrect method ": {
			method:      http.MethodGet,
			targetURL:   "/update/gauge/ram/none",
			statusCode:  http.StatusMethodNotAllowed,
			contentType: "",
		},
	}

	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			if len(tc.mockStore.args) != 0 {
				mockStore.On(tc.mockStore.methodName, tc.mockStore.args...).Return(tc.mockStore.returnArgs...)
			}

			request := resty.New().R()
			request.Method = tc.method
			u, err := url.JoinPath(srv.URL, tc.targetURL)
			assert.NoError(t, err)
			request.URL = u
			response, err := request.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.statusCode, response.StatusCode())
			assert.Equal(t, tc.contentType, response.Header().Get("Content-Type"))
			mockStore.AssertExpectations(t)

		})
	}
}

func TestUpdateJSON(t *testing.T) {
	mockStore := mocks.NewKeeper(t)

	srv := getServer(mockStore)
	defer srv.Close()

	const _addCounter = "AddCounter"
	const _setGauge = "SetGauge"
	type mockStoreRequest struct {
		methodName string
		args       []any
		returnArgs []any
	}
	tt := map[string]struct {
		method       string
		targetURL    string
		statusCode   int
		contentType  string
		mockStore    mockStoreRequest
		wantErr      error
		jsonBody     string
		jsonResponse string
	}{
		"add counter ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusOK,
			contentType: _contentTypeJSON,
			mockStore: mockStoreRequest{
				methodName: _addCounter,
				args:       []any{&metrics.Counter{Name: "1", Value: 2}},
				returnArgs: []any{&metrics.Counter{Name: "1", Value: 2}, nil},
			},
			jsonBody:     `{"id": "1", "type": "counter", "delta": 2}`,
			jsonResponse: `{"id": "1", "type": "counter", "delta": 2}`,
		},
		"zero_value ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusNotFound,
			contentType: "",
			jsonBody:    `{"id": "1", "type": "counter"}`,
		},
		"zero_metrics ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusNotFound,
			contentType: "",
			mockStore:   mockStoreRequest{},
			jsonBody:    `{"type": "counter", "delta": 2}`,
		},
		"unknown_metric_type ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusBadRequest,
			contentType: "",
			mockStore:   mockStoreRequest{},
			jsonBody:    `{"id":"er", "type": "xerp", "delta": 5}`,
		},
		"negative_float_value ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusBadRequest,
			contentType: "",
			mockStore:   mockStoreRequest{},
			jsonBody:    `{"id":"ram", "type": "counter", "delta": -2.2}`,
		},
		"simple_set_gauge ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusOK,
			contentType: _contentTypeJSON,
			mockStore: mockStoreRequest{
				methodName: _setGauge,
				args:       []any{&metrics.Gauge{Name: "ram", Value: 999.5999}},
				returnArgs: []any{&metrics.Gauge{Name: "ram", Value: 999.5999}, nil},
			},
			jsonBody: `{"id":"ram", "type": "gauge", "value": 999.5999}`,
		},
		"none_value_set_gauge ": {
			method:      http.MethodPost,
			targetURL:   "/update/",
			statusCode:  http.StatusBadRequest,
			contentType: "",
			jsonBody:    `{"id":"ram", "type": "gauge", "value": "none"}`,
		},
	}

	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			if len(tc.mockStore.args) != 0 {
				mockStore.On(tc.mockStore.methodName, tc.mockStore.args...).Return(tc.mockStore.returnArgs...)
			}

			request := resty.New().R()
			request.Method = tc.method
			u, err := url.JoinPath(srv.URL, tc.targetURL)
			assert.NoError(t, err)
			request.URL = u
			if len(tc.jsonBody) > 0 {
				request.SetHeader("Content-Type", "application/json")
				request.SetBody(tc.jsonBody)
			}
			response, err := request.Send()
			assert.NoError(t, err)
			assert.Equal(t, tc.statusCode, response.StatusCode())
			assert.Equal(t, tc.contentType, response.Header().Get("Content-Type"))
			mockStore.AssertExpectations(t)

		})
	}
}
