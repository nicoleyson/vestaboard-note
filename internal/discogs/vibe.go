package discogs

import "time"

type season int

const (
	spring season = iota
	summer
	fall
	winter
)

func seasonFor(t time.Time) season {
	switch t.Month() {
	case time.March, time.April, time.May:
		return spring
	case time.June, time.July, time.August:
		return summer
	case time.September, time.October, time.November:
		return fall
	default:
		return winter
	}
}

func vibeLabel(wmoCode int, s season) string {
	cond := conditionBucket(wmoCode)
	seasonName := [4]string{"SPRING", "SUMMER", "FALL", "WINTER"}[s]
	switch cond {
	case "clear":
		return seasonName + " SPIN"
	case "cloudy":
		return "CLOUDY " + seasonName
	case "rain":
		return "RAINY DAY PICK"
	case "snow":
		return "SNOW DAY PICK"
	case "fog":
		return "FOGGY PICK"
	case "storm":
		return "STORM PICK"
	default:
		return seasonName + " SPIN"
	}
}

func conditionBucket(wmoCode int) string {
	switch {
	case wmoCode == 0 || wmoCode == 1:
		return "clear"
	case wmoCode == 2 || wmoCode == 3:
		return "cloudy"
	case wmoCode == 45 || wmoCode == 48:
		return "fog"
	case wmoCode >= 51 && wmoCode <= 67:
		return "rain"
	case wmoCode >= 71 && wmoCode <= 77:
		return "snow"
	case wmoCode >= 80 && wmoCode <= 82:
		return "rain"
	case wmoCode >= 85 && wmoCode <= 86:
		return "snow"
	case wmoCode >= 95:
		return "storm"
	default:
		return "clear"
	}
}

var vibeStyles = map[string]map[season][]string{
	"clear": {
		spring: {"Indie Pop", "Folk", "Pop", "Psychedelic Rock", "Soft Rock"},
		summer: {"Funk", "Soul", "Reggae", "Latin", "Disco", "Afrobeat", "R&B"},
		fall:   {"Folk", "Country", "Singer/Songwriter", "Americana", "Blues"},
		winter: {"Classical", "Jazz", "Chamber Jazz", "Neo-Classical", "Piano"},
	},
	"cloudy": {
		spring: {"Indie Rock", "Dream Pop", "Shoegaze", "Alternative Rock"},
		summer: {"Indie Rock", "Alternative Rock", "Lo-Fi", "Dream Pop"},
		fall:   {"Post Bop", "Chamber Jazz", "Slowcore", "Indie Folk", "Blues"},
		winter: {"Ambient", "Post Bop", "Contemporary Jazz", "Modern Classical"},
	},
	"rain": {
		spring: {"Jazz", "Bossa Nova", "Blues", "Ambient", "Soul"},
		summer: {"Jazz", "Bossa Nova", "Blues", "Ambient", "Soul"},
		fall:   {"Jazz", "Blues", "Bossa Nova", "Ambient", "Post Bop"},
		winter: {"Jazz", "Blues", "Ambient", "Classical", "Bossa Nova"},
	},
	"fog": {
		spring: {"Ambient", "Post-Rock", "Shoegaze", "Dream Pop"},
		summer: {"Ambient", "Post-Rock", "Shoegaze", "Dream Pop"},
		fall:   {"Ambient", "Post-Rock", "Drone", "Dark Jazz"},
		winter: {"Ambient", "Drone", "Modern Classical", "Dark Jazz"},
	},
	"snow": {
		spring: {"Classical", "Folk", "Ambient", "Piano"},
		summer: {"Classical", "Folk", "Ambient", "Piano"},
		fall:   {"Classical", "Ambient", "Folk", "Piano", "Modern Classical"},
		winter: {"Classical", "Ambient", "Piano", "Modern Classical", "Neo-Classical"},
	},
	"storm": {
		spring: {"Blues", "Rock", "Electronic", "Noise Rock"},
		summer: {"Blues", "Rock", "Electronic", "Noise Rock"},
		fall:   {"Blues", "Heavy Metal", "Hard Rock", "Electronic"},
		winter: {"Blues", "Heavy Metal", "Hard Rock", "Electronic", "Industrial"},
	},
}

func stylesFor(wmoCode int, s season) []string {
	cond := conditionBucket(wmoCode)
	if styles, ok := vibeStyles[cond][s]; ok {
		return styles
	}
	return vibeStyles["clear"][s]
}

func stylesForConditionAny(wmoCode int) []string {
	cond := conditionBucket(wmoCode)
	seen := map[string]bool{}
	var all []string
	for _, seasonStyles := range vibeStyles[cond] {
		for _, s := range seasonStyles {
			if !seen[s] {
				seen[s] = true
				all = append(all, s)
			}
		}
	}
	return all
}
