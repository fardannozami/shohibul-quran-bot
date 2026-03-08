package parser

import (
	"testing"
)

func TestParse(t *testing.T) {
	p := NewReportParser()

	tests := []struct {
		message    string
		isReport   bool
		pages      int
		reportType string
	}{
		// Halaman
		{"Alhamdulillah 5 halaman", true, 5, "halaman"},
		{"alhamdulillah baca 10 hal", true, 10, "halaman"},
		{"alhamdulillah halaman 2", true, 2, "halaman"},
		{"alhamdulillah 3 hlm", true, 3, "halaman"},
		{"Alhamdulillah halaman 2-5", true, 4, "halaman"},
		{"Alhamdulillah 10 sampai 15 hal", true, 6, "halaman"},
		{"Alhamdulillah hal 1 s/d 10", true, 10, "halaman"},

		// Juz
		{"Alhamdulillah 1 juz", true, 20, "juz"},
		{"Alhamdulillah 1/2 juz", true, 10, "juz"}, 
		{"Alhamdulillah 0.5 juz", true, 10, "juz"}, 
		{"Alhamdulillah 1.5 juz", true, 30, "juz"}, 
		{"Alhamdulillah juz 1 sampai 2", true, 40, "juz"},
		{"Alhamdulillah dari juz 1 ke 3", true, 60, "juz"},
		{"alhamdulillah juz 13", true, 20, "juz"},
		{"alhamdulillah juz 13 ✅️", true, 20, "juz"},

		// Surah + Ayah
		{"Alhamdulillah surat Al-Baqarah ayat 1-30", true, 5, "surah"},
		{"alhamdulillah surat yasin", true, 6, "surah"}, // full surat: pages 440-445
		{"alhamdulillah surat kahfi ayat 1-10", true, 2, "surah"}, // pages 293-294
		{"alhamdulillah albaqoroh 1 s/d 5", true, 1, "surah"}, // page 2
		{"alhamdulillah Ali Imran ayat 10 sampai 20", true, 2, "surah"}, // pages 51-52
		{"Alhamdulillah Al 'Ankabut 46 - 79", true, 3, "surah"},         // pages 402-404 (max ayah 69)
		{"alhamdulillah al mulk", true, 3, "surah"},                     // pages 562-564
		{"alhamdulillah yasin", true, 6, "surah"},                       // pages 440-445
		{"alhamdulillah asysyuro", true, 7, "surah"},                    // pages 483-489

		// Typo / informal names
		{"alhamdulillah surat albaqoroh ayat 1-30", true, 5, "surah"},
		{"alhamdulillah yasiin 1-10", true, 1, "surah"}, // page 440

		// Default
		{"Alhamdulillah beres", true, 1, "default"},
		{"#Alhamdulillah 5 halaman", true, 5, "halaman"},
		{"#Alhamdulillaah baca 10 hal", true, 10, "halaman"},
		{"Alhamdullilah 3 hlm", true, 3, "halaman"},
		{"al hamdulillah beres", true, 1, "default"},
		{"alhamdu lillah beres", true, 1, "default"},
		{"Al-Hamdulillah beres", true, 1, "default"},

		// Not a report
		{"Bukan laporan", false, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			result := p.Parse(tt.message)

			if result.IsReport != tt.isReport {
				t.Errorf("Parse(%q).IsReport = %v; want %v", tt.message, result.IsReport, tt.isReport)
			}

			// For surah-based tests where we don't know exact pages, skip page check
			if tt.pages > 0 && result.Pages != tt.pages {
				t.Errorf("Parse(%q).Pages = %v; want %v", tt.message, result.Pages, tt.pages)
			}

			if tt.reportType != "" && result.ReportType != tt.reportType {
				t.Errorf("Parse(%q).ReportType = %q; want %q", tt.message, result.ReportType, tt.reportType)
			}
		})
	}
}

func TestFindSurahNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Al-Baqarah", 2},
		{"baqarah", 2},
		{"albaqoroh", 2},
		{"Yasin", 36},
		{"yasiin", 36},
		{"kahfi", 18},
		{"Al-Kahfi", 18},
		{"rahman", 55},
		{"ikhlas", 112},
		{"nonexistent", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			num := FindSurahNumber(tt.input)
			if num != tt.expected {
				t.Errorf("FindSurahNumber(%q) = %d; want %d", tt.input, num, tt.expected)
			}
		})
	}
}
