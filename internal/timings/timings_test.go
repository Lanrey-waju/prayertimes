package timings

import (
	"testing"
)

func TestIsPrayerTimeOver(t *testing.T) {
	if isPrayerTimeOver("05:40", "05:46") != false {
		t.Errorf("Expected false, got true")
	}
}
