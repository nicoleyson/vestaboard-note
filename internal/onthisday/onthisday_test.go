package onthisday

import (
	"testing"
	"time"
)

func TestFetch_live(t *testing.T) {
	lines, err := Fetch(time.Now())
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}
	for i, l := range lines {
		if len([]rune(l)) != 15 {
			t.Errorf("row %d: got len %d, want 15: %q", i+1, len([]rune(l)), l)
		}
	}
	t.Logf("Row1: %s", lines[0])
	t.Logf("Row2: %s", lines[1])
	t.Logf("Row3: %s", lines[2])
}
