package agent

import (
	"github.com/antoniokichaev/go-alert-me/internal/services/client/grabbers"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/mocks"
	"github.com/antoniokichaev/go-alert-me/internal/services/client/senders"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_agentBond_Run(t *testing.T) {
	mockDelivery := mocks.NewDeliveryMan(t)
	mockGrabber := mocks.NewGrabber(t)
	counter := 0
	now := func() time.Time {
		if counter > 11 {
			return time.Date(2012, time.January, 1, 1, 1, 1, 1, time.Local)
		}
		counter++
		return time.Now()
	}

	type fields struct {
		pollInterval   time.Duration
		reportInterval time.Duration
		now            func() time.Time
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
				now:            now,
				name:           "qwe",
				metricsState:   make(map[string]string),
				delivery:       mockDelivery,
				grabber:        mockGrabber,
			},
			lastState: make(map[string]string),
		},
	}
	for key, tc := range tests {
		t.Run(key, func(t *testing.T) {
			mockDelivery.EXPECT().Delivery(map[string]string{
				"ram": "5",
				"qwe": "55",
			}).Times(1).Return(nil)
			mockGrabber.EXPECT().GetSnapshot().Return(map[string]string{
				"ram": "5",
				"qwe": "55",
			}).Times(6)
			agent := &agentBond{
				pollInterval:   tc.fields.pollInterval,
				reportInterval: tc.fields.reportInterval,
				now:            tc.fields.now,
				name:           tc.fields.name,
				metricsState:   tc.fields.metricsState,
				delivery:       tc.fields.delivery,
				grabber:        tc.fields.grabber,
			}
			agent.Run()
			mockDelivery.AssertExpectations(t)
			mockGrabber.AssertExpectations(t)
			assert.EqualValues(t, tc.lastState, agent.metricsState)

		})
	}
}
