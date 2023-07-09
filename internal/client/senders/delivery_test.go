package senders

import (
	"fmt"
	"github.com/antoniokichaev/go-alert-me/internal/controller/http/v1"
	"github.com/antoniokichaev/go-alert-me/internal/entity/metrics"
	"github.com/antoniokichaev/go-alert-me/internal/usecase"
	"github.com/antoniokichaev/go-alert-me/internal/usecase/repo/mocks"
	"github.com/go-chi/chi/v5"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func getServer(mockStore *mocks.Keeper) *httptest.Server {
	getterUc := usecase.NewReceiver(mockStore)
	updaterUc := usecase.NewUpdater(mockStore)
	r := chi.NewRouter()
	v1.NewRouter(r, updaterUc, getterUc, nil, nil)
	return httptest.NewServer(r)
}

func TestLineMan_Delivery(t *testing.T) {
	mockStore := mocks.NewKeeper(t)
	srv := getServer(mockStore)
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
			mockStore.On("AddCounter", mock.Anything, &metrics.Counter{Name: "ram", Value: int64(55)}).Maybe().Return(&metrics.Counter{}, nil)
			lm := &lineMan{
				endpointRawData: tt.fields.receiver,
				httpclient:      tt.fields.httpclient,
				methodSend:      tt.fields.methodSend,
			}
			err := lm.Delivery(tt.args.data)
			tt.wantErr(t, err, fmt.Sprintf("Delivery(%v)", tt.args.data))
			mockStore.AssertExpectations(t)
		})
	}
}
