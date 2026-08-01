package uv

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		idx       float64
		wantLabel string
		wantColor int
	}{
		{0, "LOW", 66},
		{2.9, "LOW", 66},
		{3, "MODERATE", 65},
		{5.9, "MODERATE", 65},
		{6, "HIGH", 64},
		{7.9, "HIGH", 64},
		{8, "VERY HIGH", 63},
		{10.9, "VERY HIGH", 63},
		{11, "EXTREME", 68},
		{15, "EXTREME", 68},
	}
	for _, tt := range tests {
		label, color := classify(tt.idx)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%.1f) = (%q, %d), want (%q, %d)",
				tt.idx, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}
