package tearoff

import (
	"fmt"
	"strings"
	"time"
)

const (
	cRed    = 63
	cYellow = 65
	cWhite  = 69
	cols    = 15
)

var row1Tiles = [cols]int{
	cRed, cRed, cRed, cYellow, cYellow,
	cRed, cRed, cRed, cRed, cRed,
	cYellow, cYellow, cRed, cRed, cRed,
}

func Format(t time.Time) [3]string {
	var r1 strings.Builder
	for _, code := range row1Tiles {
		fmt.Fprintf(&r1, "{%d}", code)
	}

	day := t.Day()
	dayStr := fmt.Sprintf("%d", day)
	padding := cols - len(dayStr)
	leftPad := padding / 2
	rightPad := padding - leftPad
	var r2 strings.Builder
	for i := 0; i < leftPad; i++ {
		fmt.Fprintf(&r2, "{%d}", cWhite)
	}
	r2.WriteString(dayStr)
	for i := 0; i < rightPad; i++ {
		fmt.Fprintf(&r2, "{%d}", cWhite)
	}

	var r3 strings.Builder
	for i := 0; i < cols; i++ {
		fmt.Fprintf(&r3, "{%d}", cWhite)
	}

	return [3]string{r1.String(), r2.String(), r3.String()}
}
