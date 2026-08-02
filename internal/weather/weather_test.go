package weather

import (
	"testing"

	"github.com/nicoleyson/vestaboard-note/internal/geo"
)

func TestDescFromMetar(t *testing.T) {
	tests := []struct {
		cover string
		rawOb string
		want  string
	}{
		{"CLR", "KPDX 011234Z 10005KT CLR 25/10", "CLEAR"},
		{"OVC", "KPDX 011234Z 10005KT OVC010", "OVERCAST"},
		{"BKN", "KPDX 011234Z BKN020", "CLOUDY"},
		{"SCT", "KPDX 011234Z SCT030", "PARTLY CLOUDY"},
		{"FEW", "KPDX 011234Z FEW040", "MOSTLY CLEAR"},
		{"OVC", "KPDX 011234Z 10005KT OVC010 RA", "RAIN"},
		{"OVC", "KPDX 011234Z 10005KT OVC010 SN", "SNOW"},
		{"OVC", "KPDX 011234Z FG OVC002", "FOG"},
		{"OVC", "KPDX 011234Z TSRA OVC010", "THUNDERSTORM"},
		{"OVC", "KPDX 011234Z FZRA OVC010", "FREEZING RAIN"},
		{"OVC", "KPDX 011234Z BLSN OVC010", "BLIZZARD"},
		{"OVC", "KPDX 011234Z SHRA OVC015", "SHOWERS"},
		{"OVC", "KPDX 011234Z RASN OVC010", "RAIN AND SNOW"},
	}
	for _, tt := range tests {
		got := descFromMetar(tt.cover, tt.rawOb)
		if got != tt.want {
			t.Errorf("descFromMetar(%q, %q) = %q, want %q", tt.cover, tt.rawOb, got, tt.want)
		}
	}
}

func TestColorForDesc(t *testing.T) {
	tests := []struct {
		desc string
		want int
	}{
		{"CLEAR", 65},
		{"MOSTLY CLEAR", 65},
		{"CLOUDY", 69},
		{"OVERCAST", 69},
		{"RAIN", 67},
		{"SHOWERS", 67},
		{"DRIZZLE", 67},
		{"SNOW", 69},
		{"FREEZING RAIN", 69},
		{"BLIZZARD", 69},
		{"FOG", 69},
		{"THUNDERSTORM", 68},
	}
	for _, tt := range tests {
		got := colorForDesc(tt.desc)
		if got != tt.want {
			t.Errorf("colorForDesc(%q) = %d, want %d", tt.desc, got, tt.want)
		}
	}
}

func TestIsFahrenheitCountry(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lon  float64
		want bool
	}{
		{"Portland OR", 45.5, -122.6, true},
		{"New York NY", 40.7, -74.0, true},
		{"Miami FL", 25.8, -80.2, true},
		{"Honolulu HI", 21.3, -157.8, true},
		{"Anchorage AK", 61.2, -149.9, true},
		{"London UK", 51.5, -0.1, false},
		{"Tokyo JP", 35.7, 139.7, false},
		{"Paris FR", 48.9, 2.3, false},
		{"Sydney AU", -33.9, 151.2, false},
	}
	for _, tt := range tests {
		got := geo.IsFahrenheitCountry(tt.lat, tt.lon)
		if got != tt.want {
			t.Errorf("IsFahrenheitCountry(%s: %.1f, %.1f) = %v, want %v", tt.name, tt.lat, tt.lon, got, tt.want)
		}
	}
}

func TestDist(t *testing.T) {
	d := dist(0, 0, 3, 4)
	if d != 25 {
		t.Errorf("dist(0,0,3,4) = %f, want 25", d)
	}
	if dist(1, 1, 1, 1) != 0 {
		t.Errorf("dist of same point should be 0")
	}
}
