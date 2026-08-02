package geo

import "testing"

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
		{"Anchorage AK", 61.2, -149.9, true},
		{"Honolulu HI", 21.3, -157.8, true},
		{"San Juan PR", 18.4, -66.1, true},
		{"St. Thomas USVI", 18.3, -64.9, true},
		{"Guam", 13.5, 144.8, true},
		{"American Samoa", -14.3, -170.7, true},
		{"Monrovia Liberia", 6.3, -10.8, true},
		{"London UK", 51.5, -0.1, false},
		{"Tokyo JP", 35.7, 139.7, false},
		{"Paris FR", 48.9, 2.3, false},
		{"Sydney AU", -33.9, 151.2, false},
		{"Mexico City MX", 19.4, -99.1, false},
		{"North Pole", 90.0, 0.0, false},
	}
	for _, tt := range tests {
		got := IsFahrenheitCountry(tt.lat, tt.lon)
		if got != tt.want {
			t.Errorf("IsFahrenheitCountry(%s: %.1f, %.1f) = %v, want %v",
				tt.name, tt.lat, tt.lon, got, tt.want)
		}
	}
}

