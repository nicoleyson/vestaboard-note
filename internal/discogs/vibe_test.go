package discogs

import (
	"testing"
	"time"
)

func TestSeasonFor(t *testing.T) {
	tests := []struct {
		month time.Month
		want  season
	}{
		{time.March, spring},
		{time.April, spring},
		{time.May, spring},
		{time.June, summer},
		{time.July, summer},
		{time.August, summer},
		{time.September, fall},
		{time.October, fall},
		{time.November, fall},
		{time.December, winter},
		{time.January, winter},
		{time.February, winter},
	}
	for _, tt := range tests {
		got := seasonFor(time.Date(2024, tt.month, 1, 0, 0, 0, 0, time.UTC))
		if got != tt.want {
			t.Errorf("seasonFor(%v) = %d, want %d", tt.month, got, tt.want)
		}
	}
}

func TestTimeSlotFor(t *testing.T) {
	tests := []struct {
		hour int
		want timeSlot
	}{
		{0, lateNight},
		{4, lateNight},
		{5, earlyMorning},
		{8, earlyMorning},
		{9, morning},
		{11, morning},
		{12, afternoon},
		{16, afternoon},
		{17, evening},
		{19, evening},
		{20, night},
		{23, night},
	}
	for _, tt := range tests {
		base := time.Date(2024, time.January, 1, tt.hour, 0, 0, 0, time.UTC)
		got := timeSlotFor(base)
		if got != tt.want {
			t.Errorf("timeSlotFor(hour=%d) = %d, want %d", tt.hour, got, tt.want)
		}
	}
}

func TestConditionBucket(t *testing.T) {
	tests := []struct {
		wmo  int
		want string
	}{
		{0, "clear"},
		{1, "clear"},
		{2, "cloudy"},
		{3, "cloudy"},
		{45, "fog"},
		{48, "fog"},
		{51, "rain"},
		{65, "rain"},
		{67, "rain"},
		{71, "snow"},
		{77, "snow"},
		{80, "rain"},
		{82, "rain"},
		{85, "snow"},
		{86, "snow"},
		{95, "storm"},
		{99, "storm"},
		{10, "clear"},
	}
	for _, tt := range tests {
		got := conditionBucket(tt.wmo)
		if got != tt.want {
			t.Errorf("conditionBucket(%d) = %q, want %q", tt.wmo, got, tt.want)
		}
	}
}

func TestVibeLabelNotEmpty(t *testing.T) {
	seasons := []season{spring, summer, fall, winter}
	slots := []timeSlot{earlyMorning, morning, afternoon, evening, night, lateNight}
	wmos := []int{0, 2, 45, 55, 73, 81, 95}
	for _, s := range seasons {
		for _, sl := range slots {
			for _, w := range wmos {
				label := vibeLabel(w, s, sl)
				if label == "" {
					t.Errorf("vibeLabel(wmo=%d, season=%d, slot=%d) returned empty string", w, s, sl)
				}
			}
		}
	}
}

func TestTitleMatchesSlot(t *testing.T) {
	if !titleMatchesSlot("Morning Light", morning) {
		t.Error("titleMatchesSlot: 'Morning Light' should match morning slot")
	}
	if !titleMatchesSlot("Midnight Rambler", lateNight) {
		t.Error("titleMatchesSlot: 'Midnight Rambler' should match lateNight slot")
	}
	if titleMatchesSlot("Blue Train", morning) {
		t.Error("titleMatchesSlot: 'Blue Train' should not match morning slot")
	}
}
