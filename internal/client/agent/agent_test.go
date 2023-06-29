package agent

import (
	"github.com/antoniokichaev/go-alert-me/internal/client/grabbers"
	grabber_mocks "github.com/antoniokichaev/go-alert-me/internal/client/grabbers/mocks"
	"github.com/antoniokichaev/go-alert-me/internal/client/senders"
	sender_mocks "github.com/antoniokichaev/go-alert-me/internal/client/senders/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_agentBond_Run(t *testing.T) {
	mockDelivery := sender_mocks.NewDeliveryMan(t)
	mockGrabber := grabber_mocks.NewGrabber(t)

	type fields struct {
		pollInterval   time.Duration
		reportInterval time.Duration
		name           string
		metricsState   map[string]string
		delivery       senders.DeliveryMan
		grabber        grabbers.Grabber
	}
	tests := map[string]struct {
		fields    fields
		lastState map[string]string
	}{
		"check_call_delivery": {
			fields: fields{
				pollInterval:   time.Second * 1,
				reportInterval: time.Second * 5,
				name:           "qwe",
				metricsState:   make(map[string]string),
				grabber:        mockGrabber,
				delivery:       mockDelivery,
			},
			lastState: make(map[string]string),
		},
	}
	//
	//minCallsGetSnapshot := 10
	//minCallsDeliveryBody := 2
	for key, tc := range tests {
		t.Run(key, func(t *testing.T) {
			//mockDelivery.On("DeliveryBody", mock.Anything).Return(nil).Maybe()
			mockDelivery.EXPECT().DeliveryBody(mock.Anything).Return(nil).Maybe()
			mockGrabber.EXPECT().GetSnapshot().Return(map[string]string{
				"counter/test": "5",
				"gauge/age":    "55",
			}).Maybe()

			stop := make(chan struct{})
			agent := &agentBond{
				pollInterval:   tc.fields.pollInterval,
				reportInterval: tc.fields.reportInterval,
				name:           tc.fields.name,
				metricsState:   tc.fields.metricsState,
				delivery:       tc.fields.delivery,
				grabber:        tc.fields.grabber,
				notify:         stop,
			}
			tick := time.NewTicker(time.Second * 15)
			go func() {
				<-tick.C
				stop <- struct{}{}
			}()
			agent.Run()
			mockDelivery.AssertExpectations(t)
			mockGrabber.AssertExpectations(t)
			mockDelivery.AssertCalled(t, "DeliveryBody", mock.Anything)
			mockGrabber.AssertCalled(t, "GetSnapshot")

		})
	}
}
