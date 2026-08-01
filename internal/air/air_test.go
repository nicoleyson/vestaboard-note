package air

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		aqi       int
		wantLabel string
		wantColor int
	}{
		{0, "GOOD", 66},
		{50, "GOOD", 66},
		{51, "MODERATE", 65},
		{100, "MODERATE", 65},
		{101, "UNHEALTHY+", 64},
		{150, "UNHEALTHY+", 64},
		{151, "UNHEALTHY", 63},
		{200, "UNHEALTHY", 63},
		{201, "VERY UNHLTHY", 68},
		{300, "VERY UNHLTHY", 68},
		{301, "HAZARDOUS", 70},
		{500, "HAZARDOUS", 70},
	}
	for _, tt := range tests {
		label, color := classify(tt.aqi)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%d) = (%q, %d), want (%q, %d)",
				tt.aqi, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}
