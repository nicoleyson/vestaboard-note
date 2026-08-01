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

func TestClassifyTrivial(t *testing.T) {
	trivialCases := []int{0, 25, 50}
	for _, v := range trivialCases {
		label, _ := classify(v)
		if label != "GOOD" {
			t.Errorf("classify(%d) = %q, expected GOOD (trivial)", v, label)
		}
	}
	nonTrivialCases := []int{51, 100, 200}
	for _, v := range nonTrivialCases {
		label, _ := classify(v)
		if label == "GOOD" {
			t.Errorf("classify(%d) = GOOD, expected non-trivial", v)
		}
	}
}
