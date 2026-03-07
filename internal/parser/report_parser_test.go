package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	p := NewReportParser()

	tests := []struct {
		message string
		isReport bool
		pages int
	}{
		{"Alhamdulillah 5 halaman", true, 5},
		{"alhamdulillah baca 10 hal", true, 10},
		{"Alhamdulillah 1 juz", true, 20},
		{"Bukan laporan", false, 0},
		{"Alhamdulillah beres", true, 1}, // default 1
		{"alhamdulillah halaman 2", true, 2},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			isR, pgs := p.Parse(tt.message)
			if isR != tt.isReport || pgs != tt.pages {
				t.Errorf("Parse(%q) = %v, %v; want %v, %v", tt.message, isR, pgs, tt.isReport, tt.pages)
			}
		})
	}
}
