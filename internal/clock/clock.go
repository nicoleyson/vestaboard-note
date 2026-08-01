package clock

import (
	"time"

	"github.com/nicoleyson/vestaboard-note/internal/layout"
)

func Format(t time.Time) [3]string {
	return [3]string{
		layout.Center(t.Format("MON JAN 2"), layout.Cols),
		layout.Center(t.Format("3:04 PM"), layout.Cols),
		layout.Center(t.Format("2006"), layout.Cols),
	}
}
