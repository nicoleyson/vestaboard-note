package rain

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		prob      int
		mm        float64
		wantLabel string
		wantColor int
	}{
		{0, 0, "NONE", 65},
		{20, 0, "NONE", 65},
		{30, 0, "POSSIBLE", 66},
		{59, 0, "POSSIBLE", 66},
		{60, 0, "LIKELY", 67},
		{100, 0, "LIKELY", 67},
		{80, 0.1, "LIGHT", 67},
		{80, 1.0, "MODERATE", 67},
		{80, 4.0, "HEAVY", 67},
		{80, 10.0, "HEAVY", 67},
	}
	for _, tt := range tests {
		label, color := classify(tt.prob, tt.mm)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(prob=%d, mm=%.1f) = (%q, %d), want (%q, %d)",
				tt.prob, tt.mm, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}
