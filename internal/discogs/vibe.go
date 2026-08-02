package discogs

import (
	"strings"
	"time"
)

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

type timeSlot int

const (
	earlyMorning timeSlot = iota
	morning
	afternoon
	evening
	night
	lateNight
)

func timeSlotFor(t time.Time) timeSlot {
	h := t.Hour()
	switch {
	case h >= 5 && h < 9:
		return earlyMorning
	case h >= 9 && h < 12:
		return morning
	case h >= 12 && h < 17:
		return afternoon
	case h >= 17 && h < 20:
		return evening
	case h >= 20:
		return night
	default:
		return lateNight
	}
}

var timeSlotStyles = map[timeSlot][]string{
	earlyMorning: {"Ambient", "Acoustic", "Folk", "Lo-Fi", "Contemporary Jazz", "Piano", "Bossa Nova", "New Age"},
	morning:      {"Pop", "Soul", "Indie Pop", "Funk", "R&B", "Folk", "Soft Rock"},
	afternoon:    {"Rock", "Funk", "Soul", "Latin", "Indie Rock", "Alternative Rock", "Pop"},
	evening:      {"Singer/Songwriter", "Jazz", "R&B", "Soul", "Indie Folk", "Bossa Nova", "Contemporary Jazz"},
	night:        {"Electronic", "Post Bop", "Slowcore", "Dream Pop", "Ambient", "Neo-Soul", "Dark Jazz"},
	lateNight:    {"Ambient", "Drone", "Post Bop", "Electronic", "Dark Jazz", "Contemporary Jazz", "Lo-Fi"},
}

var timeSlotKeywords = map[timeSlot][]string{
	earlyMorning: {"5AM", "6AM", "7AM", "8AM", "MORNING", "DAWN", "SUNRISE", "EARLY", "RISE"},
	morning:      {"MORNING", "COFFEE", "BREAKFAST", "AM", "WAKE"},
	afternoon:    {"AFTERNOON", "MIDDAY", "LUNCH", "SUNDAY", "LAZY"},
	evening:      {"EVENING", "SUNSET", "DUSK", "GOLDEN HOUR", "TWILIGHT"},
	night:        {"NIGHT", "MIDNIGHT", "AFTER DARK", "LATE NIGHT", "LIGHTS OUT"},
	lateNight:    {"MIDNIGHT", "4AM", "3AM", "2AM", "1AM", "LATE NIGHT", "INSOMNIA", "SLEEPLESS", "INSOMNIAC", "AFTER HOURS"},
}

var timeSlotNames = map[timeSlot]string{
	earlyMorning: "EARLY MORNING",
	morning:      "MORNING",
	afternoon:    "AFTERNOON",
	evening:      "EVENING",
	night:        "NIGHT",
	lateNight:    "LATE NIGHT",
}

func titleMatchesSlot(title string, slot timeSlot) bool {
	upper := strings.ToUpper(title)
	for _, kw := range timeSlotKeywords[slot] {
		if strings.Contains(upper, kw) {
			return true
		}
	}
	return false
}

func vibeLabel(wmoCode int, s season, slot timeSlot) string {
	cond := conditionBucket(wmoCode)
	slotName := timeSlotNames[slot]
	switch cond {
	case "clear":
		seasonName := [4]string{"SPRING", "SUMMER", "FALL", "WINTER"}[s]
		if slot == lateNight || slot == night {
			return slotName + " SPIN"
		}
		return seasonName + " " + slotName
	case "rain":
		return "RAINY " + slotName
	case "snow":
		return "SNOWY " + slotName
	case "fog":
		return "FOGGY " + slotName
	case "storm":
		return "STORM PICK"
	case "cloudy":
		if slot == lateNight || slot == night {
			return slotName + " SPIN"
		}
		return "CLOUDY " + slotName
	default:
		return slotName + " SPIN"
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
