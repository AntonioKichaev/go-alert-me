package senders

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/services/server/handlers/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/storages/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestLineMan_Delivery(t *testing.T) {
	mockStore := mocks.NewMetricRepository(t)
	handlerMetrics := metrics.NewHandlerMetrics(mockStore)
	r := chi.NewRouter()
	handlerMetrics.Register(r)
	srv := httptest.NewServer(r)
	targetURL, err := url.JoinPath(srv.URL, "/update")
	client := resty.NewWithClient(srv.Client())
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		receiver   string
		httpclient *resty.Client
		methodSend string
	}
	type args struct {
		data map[string]string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "simple send",
			fields: fields{
				receiver:   targetURL,
				httpclient: client,
				methodSend: http.MethodPost,
			},
			args: args{
				data: map[string]string{"counter/ram": "55"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "not found",
			fields: fields{
				receiver:   targetURL,
				httpclient: client,
				methodSend: http.MethodPost,
			},
			args: args{
				data: map[string]string{"ram": "test"},
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				assert.ErrorIs(t, err, ErrorStatusCode, i...)
				return false
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore.On("AddCounter", "ram", int64(55)).Maybe()
			lm := &lineMan{
				receiver:   tt.fields.receiver,
				httpclient: tt.fields.httpclient,
				methodSend: tt.fields.methodSend,
			}
			tt.wantErr(t, lm.Delivery(tt.args.data), fmt.Sprintf("Delivery(%v)", tt.args.data))
		})
	}
}
