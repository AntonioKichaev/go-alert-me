package metrics

import (
	"github.com/antoniokichaev/go-alert-me/internal/storages/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestUpdateMetrics(t *testing.T) {
	mockStore := mocks.NewMetricRepository(t)
	handlerMetrics := newHandlerMetrics(mockStore)
	r := chi.NewRouter()
	handlerMetrics.Register(r)

	const _addCounter = "AddCounter"
	const _setGauge = "SetGauge"
	type mockStoreRequest struct {
		methodName string
		args       []any
	}
	tt := map[string]struct {
		method      string
		targetURL   string
		statusCode  int
		contentType string
		mockStore   mockStoreRequest
	}{
		"add counter ": {
			method:    http.MethodPost,
			targetURL: "/update/counter/1/2",

			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore:   mockStoreRequest{methodName: _addCounter, args: []any{"1", int64(2)}},
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
			statusCode:  http.StatusNotFound,
			contentType: _contentTypeText,
			mockStore:   mockStoreRequest{methodName: _addCounter},
		},
		"negative_value ": {
			method:      http.MethodPost,
			targetURL:   "/update/counter/ram/-5",
			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore:   mockStoreRequest{methodName: _addCounter, args: []any{"ram", int64(-5)}},
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
			mockStore:   mockStoreRequest{methodName: _setGauge, args: []any{"ram", 999.5999}},
		},
		"none_value_set_gauge ": {
			method:      http.MethodPost,
			targetURL:   "/update/gauge/ram/none",
			statusCode:  http.StatusBadRequest,
			contentType: "",
		},
	}

	srv := httptest.NewServer(r)
	defer srv.Close()

	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			if len(tc.mockStore.args) != 0 {

				mockStore.On(tc.mockStore.methodName, tc.mockStore.args...)
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

func TestHandlerMetric(t *testing.T) {
	mockStore := mocks.NewMetricRepository(t)
	handlerCounter := newHandlerMetrics(mockStore)
	handler := http.HandlerFunc(handlerCounter.updateMetrics)
	srv := httptest.NewServer(handler)
	defer srv.Close()
}
