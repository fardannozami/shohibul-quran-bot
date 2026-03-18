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
		{"alhamdulillah halaman 2", true, 1, "halaman"}, // 'halaman 2' means page 2 (1 page)
		{"alhamdulillah 3 hlm", true, 3, "halaman"},
		{"Alhamdulillah halaman 2-5", true, 4, "halaman"},
		{"Alhamdulillah 10 sampai 15 hal", true, 6, "halaman"},
		{"Alhamdulillah hal 1 s/d 10", true, 10, "halaman"},
		{"alhamdulillah 4 halaman hari ini", true, 4, "halaman"},
		{"alhamdulillah 4 halaman tadi malam", true, 4, "halaman"},

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
		{"alhamdulillah surat yasin", true, 6, "surah"},                 // full surat: pages 440-445
		{"alhamdulillah surat kahfi ayat 1-10", true, 2, "surah"},       // pages 293-294
		{"alhamdulillah albaqoroh 1 s/d 5", true, 1, "surah"},           // page 2
		{"alhamdulillah Ali Imran ayat 10 sampai 20", true, 2, "surah"}, // pages 51-52
		{"Alhamdulillah Al 'Ankabut 46 - 79", true, 3, "surah"},         // pages 402-404 (max ayah 69)
		{"alhamdulillah al mulk", true, 3, "surah"},                     // pages 562-564
		{"alhamdulillah yasin", true, 6, "surah"},                       // pages 440-445
		{"alhamdulillah asysyuro", true, 7, "surah"},                    // pages 483-489
		{"alhamdulillah asy-syura", true, 7, "surah"},                   // pages 483-489
		{"alhamdulillah asy-syu'ara'", true, 10, "surah"},              // pages 367-376 (max ayah 227)

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

		// Surah Range
		{"Alhamdulillah al ashr - an nas", true, 4, "surah"},
		{"Alhamdulillah Al-Fatihah sampai Al-Baqarah", true, 49, "surah"},

		{"Alhamdulillah alfatihah sampai albaqorah ayat 100", true, 15, "surah"},

		// Multi-line Juz list (take latest)
		{`#Alhamdulillah hari Ramadhan
1. Juz 21 ✅ Juz 22✅
2. Juz 23✅
3. Juz 24✅
4. Juz 25✅
5. Juz 26 ✅
6. Juz 27✅
7. Juz 28✅
8. Juz 29✅
9. Juz 30✅
10. Juz 1 ✅
11. Juz 2 ✅
12. Juz 3✅JUZ 4 ✅
13. Juz 5✅ Juz 6✅
14. Juz 7✅ Juz 8✅
15. Juz 9✅ Juz 10✅
16. Juz 11✅+ Al-Kahf✅ 
17. Juz 12✅
18. Juz 13✅
19. Juz 14✅
20. Juz 15,16,17,18✅
21. Juz 19, 20✅
22. Juz 21✅
23. Juz 22✅
24. Juz 23,24✅
25. Juz 25✅
26. Juz 26-30✅
27. Juz 1✅
28. Juz 2✅
29. Juz 3✅`, true, 20, "juz"},

		// Not a report
		{"Bukan laporan", false, 0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.message, func(t *testing.T) {
			results := p.Parse(tt.message)

			if !tt.isReport {
				if len(results) > 0 {
					t.Errorf("Parse(%q) got %d results, want 0", tt.message, len(results))
				}
				return
			}

			if len(results) == 0 {
				t.Errorf("Parse(%q) got no results, want report", tt.message)
				return
			}

			// Sum total pages for verification
			totalPages := 0
			for _, r := range results {
				totalPages += r.Pages
			}

			if tt.pages > 0 && totalPages != tt.pages {
				t.Errorf("Parse(%q).TotalPages = %v; want %v", tt.message, totalPages, tt.pages)
			}

			if tt.reportType != "" && results[0].ReportType != tt.reportType {
				t.Errorf("Parse(%q).ReportType = %q; want %q", tt.message, results[0].ReportType, tt.reportType)
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
		{"ankaboot", 29}, // fuzzy Al-'Ankabut
		{"fatiha", 1},    // fuzzy Al-Fatihah
		{"rohman", 55},   // will match Ar-Rahman (score should be high enough)
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
