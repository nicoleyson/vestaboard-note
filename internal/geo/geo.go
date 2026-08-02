package geo

// IsFahrenheitCountry returns true for countries that use Fahrenheit:
// United States (continental, Alaska, Hawaii, territories) and Liberia.
func IsFahrenheitCountry(lat, lon float64) bool {
	// Continental US
	if lat >= 24 && lat <= 49.5 && lon >= -125 && lon <= -66 {
		return true
	}
	// Alaska
	if lat >= 54 && lat <= 72 && lon >= -168 && lon <= -130 {
		return true
	}
	// Hawaii
	if lat >= 18 && lat <= 23 && lon >= -161 && lon <= -154 {
		return true
	}
	// Puerto Rico + US Virgin Islands
	if lat >= 17 && lat <= 18.5 && lon >= -68 && lon <= -64 {
		return true
	}
	// Liberia
	if lat >= 4 && lat <= 8.5 && lon >= -11.5 && lon <= -7.5 {
		return true
	}
	return false
}
