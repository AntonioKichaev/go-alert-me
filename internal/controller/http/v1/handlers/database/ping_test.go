package database_test

import (
	"errors"
	v1 "github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1/handlers/database/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHadlerDB_GetPing(t *testing.T) {
	db := mocks.NewStorageStatus(t)
	r := chi.NewRouter()
	v1.NewRouter(r, nil, nil, db, nil)

	srv := httptest.NewServer(r)
	defer srv.Close()

	tests := map[string]struct {
		returnValue error
		wantStatus  int
	}{
		"1": {
			wantStatus:  http.StatusInternalServerError,
			returnValue: errors.New("asdadad"),
		},
		"2": {
			returnValue: nil,
			wantStatus:  http.StatusOK,
		},
		"3": {
			returnValue: nil,
			wantStatus:  http.StatusOK,
		},
	}
	for key, tt := range tests {
		t.Run(key, func(t *testing.T) {
			db.EXPECT().Ping().Return(tt.returnValue).Once()
			request := resty.New().R()
			request.Method = http.MethodGet
			u, err := url.JoinPath(srv.URL, "/ping")
			assert.NoError(t, err)
			request.URL = u

			response, err := request.Send()
			assert.NoError(t, err)

			assert.Equal(t, response.StatusCode(), tt.wantStatus)

		})
	}
}
