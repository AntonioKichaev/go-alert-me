package metrics_test

import (
	"errors"
	v1 "github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/entity"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/antoniokichaev/go-alert-me/internal/usecase/repo/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

const _contentTypeText = "text/plain; charset=utf-8"

func getServer(mockStore *mocks.Keeper) *httptest.Server {
	getterUc := usecase.NewReceiver(mockStore)
	updaterUc := usecase.NewUpdater(mockStore)
	r := chi.NewRouter()
	v1.NewRouter(r, updaterUc, getterUc)
	return httptest.NewServer(r)
}

func TestGetMetrics(t *testing.T) {
	mockStore := mocks.NewKeeper(t)
	srv := getServer(mockStore)
	defer srv.Close()

	const _getGauge = "GetGauge"
	const _getCounter = "GetCounter"
	type mockStoreRequest struct {
		methodName  string
		args        []any
		returnValue []any
	}

	tt := map[string]struct {
		method      string
		targetURL   string
		statusCode  int
		contentType string
		wantErr     bool
		mockStore   mockStoreRequest
	}{
		"exist counter": {
			method:      http.MethodGet,
			targetURL:   "/value/counter/my",
			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore: mockStoreRequest{
				methodName:  _getCounter,
				args:        []any{"my"},
				returnValue: []any{&entity.Counter{Name: "my", Value: 5}, nil},
			},
			wantErr: false,
		},
		"exist gauge": {
			method:      http.MethodGet,
			targetURL:   "/value/gauge/my",
			statusCode:  http.StatusOK,
			contentType: _contentTypeText,
			mockStore: mockStoreRequest{
				methodName:  _getGauge,
				args:        []any{"my"},
				returnValue: []any{&entity.Gauge{Name: "my", Value: 5}, nil},
			},
			wantErr: false,
		},
		"doesn't exist ": {
			method:      http.MethodGet,
			targetURL:   "/value/counter/unk",
			statusCode:  http.StatusNotFound,
			contentType: "",
			wantErr:     true,
			mockStore: mockStoreRequest{
				methodName:  _getCounter,
				args:        []any{"unk"},
				returnValue: []any{nil, errors.New("NotFound")},
			},
		},
	}

	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			if len(tc.mockStore.args) != 0 || tc.mockStore.methodName != "" {
				mockStore.On(tc.mockStore.methodName, tc.mockStore.args...).Return(tc.mockStore.returnValue...)
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

func TestGetAllMetrics(t *testing.T) {
	mockStore := mocks.NewKeeper(t)
	srv := getServer(mockStore)
	defer srv.Close()

	tt := map[string]struct {
		returnStore map[string]string
		want        string
		wantErr     error
	}{
		"get all metrics": {
			returnStore: map[string]string{
				"lex":  "26",
				"test": "25",
			},
			want:    "test 25\nlex 26\n",
			wantErr: nil,
		},
		"get nothing metrics": {
			returnStore: nil,
			want:        "",
			wantErr:     nil,
		},
	}
	for key, tc := range tt {
		t.Run(key, func(t *testing.T) {
			mockStore.EXPECT().GetMetrics().Return(tc.returnStore, tc.wantErr)
			request := resty.New().R()
			request.Method = http.MethodGet
			request.URL = srv.URL
			response, err := request.Send()
			got := strings.Split(string(response.Body()), "\n")
			want := strings.Split(tc.want, "\n")
			assert.NoError(t, err)
			assert.ObjectsAreEqualValues(want, got)

		})
	}

}