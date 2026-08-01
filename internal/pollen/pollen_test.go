package pollen

import (
	"testing"
)

func TestClassify(t *testing.T) {
	tests := []struct {
		val       float64
		wantLabel string
		wantColor int
	}{
		{0, "LOW", 66},
		{9.9, "LOW", 66},
		{10, "MODERATE", 65},
		{49.9, "MODERATE", 65},
		{50, "HIGH", 64},
		{199.9, "HIGH", 64},
		{200, "VERY HIGH", 63},
		{500, "VERY HIGH", 63},
	}
	for _, tt := range tests {
		label, color := classify(tt.val)
		if label != tt.wantLabel || color != tt.wantColor {
			t.Errorf("classify(%.1f) = (%q, %d), want (%q, %d)",
				tt.val, label, color, tt.wantLabel, tt.wantColor)
		}
	}
}

func TestDominantType(t *testing.T) {
	tests := []struct {
		grass, tree, weed float64
		wantName          string
	}{
		{10, 5, 3, "GRASS"},
		{5, 10, 3, "TREE"},
		{3, 5, 10, "WEED"},
		{10, 10, 5, "GRASS"},
		{5, 10, 10, "TREE"},
		{0, 0, 0, "GRASS"},
	}
	for _, tt := range tests {
		name, _ := dominantType(tt.grass, tt.tree, tt.weed)
		if name != tt.wantName {
			t.Errorf("dominantType(%.0f, %.0f, %.0f) = %q, want %q",
				tt.grass, tt.tree, tt.weed, name, tt.wantName)
		}
	}
}

func TestClassifyTrivial(t *testing.T) {
	trivialCases := []float64{0, 5, 9.9}
	for _, v := range trivialCases {
		label, _ := classify(v)
		if label != "LOW" {
			t.Errorf("classify(%.1f) = %q, expected LOW (trivial)", v, label)
		}
	}
	nonTrivialCases := []float64{10, 50, 200}
	for _, v := range nonTrivialCases {
		label, _ := classify(v)
		if label == "LOW" {
			t.Errorf("classify(%.1f) = LOW, expected non-trivial", v)
		}
	}
}
