package scheduler

import (
	"strings"
	"testing"
	"time"
)

func TestRamadhanLine(t *testing.T) {
	start := time.Date(2027, 2, 8, 0, 0, 0, 0, time.Local)
	at := func(y int, m time.Month, d int) time.Time { return time.Date(y, m, d, 7, 0, 0, 0, time.Local) }

	tests := []struct {
		name string
		now  time.Time
		want string
	}{
		{"far away", at(2026, 10, 1), ""},
		{"countdown", at(2026, 11, 16), "84 hari lagi"},
		{"tomorrow", at(2027, 2, 7), "besok"},
		{"first day", at(2027, 2, 8), "hari pertama"},
		{"day 10", at(2027, 2, 17), "hari ke-10"},
		{"day 30", at(2027, 3, 9), "hari ke-30"},
		{"after", at(2027, 3, 10), ""},
	}
	for _, tt := range tests {
		got := ramadhanLine(start, tt.now)
		if tt.want == "" && got != "" {
			t.Errorf("%s: want empty, got %q", tt.name, got)
		}
		if tt.want != "" && !strings.Contains(got, tt.want) {
			t.Errorf("%s: want %q in %q", tt.name, tt.want, got)
		}
	}
}
