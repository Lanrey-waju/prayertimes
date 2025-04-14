package timings

import (
	"testing"
)

func TestIsPrayerTimeOver(t *testing.T) {
	tests := map[string]struct {
		currentTime string
		prayerTime  string
		want        bool
	}{
		"current time before prayer time": {currentTime: "05:06", prayerTime: "05:40", want: false},
		"current time after prayer time":  {currentTime: "13:00", prayerTime: "12:58", want: true},
		"current time equals prayer time": {currentTime: "13:00", prayerTime: "13:00", want: false},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := isPrayerTimeOver(tc.prayerTime, tc.currentTime)
			if tc.want != got {
				t.Fatalf("expected %v, got %v", tc.want, got)
			}
		})
	}
}
