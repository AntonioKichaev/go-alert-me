package client

import (
	"github.com/antoniokichaev/go-alert-me/internal/services/client/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRacoon_GetSnapshot(t *testing.T) {
	rnd := mocks.NewRandom(t)

	tests := []struct {
		name         string
		randomVal    int
		pollCount    int
		wantSnapShot map[string]string
	}{
		{
			name:      "common behave",
			randomVal: 324,
			pollCount: 55,
			wantSnapShot: map[string]string{
				"RandomValue": "324",
				"PollCount":   "55",
			},
		},
		{
			name:         "zero_calls",
			randomVal:    324,
			pollCount:    0,
			wantSnapShot: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pollCount != 0 {
				rnd.EXPECT().Int().Return(tt.randomVal)
			}
			racoon := Racoon{
				PollCount: 0,
				random:    rnd,
			}
			gotSnap := map[string]string{}
			for i := 0; i < tt.pollCount; i++ {
				gotSnap = racoon.GetSnapshot()
			}

			for key, val := range tt.wantSnapShot {
				assert.Contains(t, gotSnap, key, "GetSnapshot()")
				assert.Equal(t, val, gotSnap[key])
			}
			assert.Equal(t, tt.pollCount, racoon.PollCount)
			rnd.AssertExpectations(t)
		})
	}
}
