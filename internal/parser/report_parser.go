package parser

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseResult holds the result of parsing a message
type ParseResult struct {
	IsReport   bool
	Pages      int
	SurahName  string // e.g. "Al-Baqarah" (empty if not surah-based)
	StartAyah  int    // e.g. 1 (0 if not surah-based)
	EndAyah    int    // e.g. 30 (0 if not surah-based)
	ReportType string // "halaman", "juz", "surah", or "default"
}

// ReportParser handles parsing of chat messages to detect reading activity
type ReportParser struct{}

func NewReportParser() *ReportParser {
	return &ReportParser{}
}

// Parse determines if the message contains a valid report and returns a ParseResult
func (p *ReportParser) Parse(message string) ParseResult {
	lowerMsg := strings.ToLower(message)

	// Must contain alhamdulillah (flexible match)
	alhamdulillahRegex := regexp.MustCompile(`(?i)#?al[ -]?hamdu?[ -]?l+il+a+h`)
	if !alhamdulillahRegex.MatchString(message) {
		return ParseResult{}
	}

	// 1. Try halaman
	if pages := p.extractPages(lowerMsg); pages > 0 {
		return ParseResult{IsReport: true, Pages: pages, ReportType: "halaman"}
	}

	// 2. Try juz
	if pages := p.extractJuz(lowerMsg); pages > 0 {
		return ParseResult{IsReport: true, Pages: pages, ReportType: "juz"}
	}

	// 3. Try surah + ayah
	if result := p.extractSurahAyah(message); result.IsReport {
		return result
	}

	// 4. Default: assume 1 page if alhamdulillah detected
	return ParseResult{IsReport: true, Pages: 1, ReportType: "default"}
}

// ParseCompat is a backward-compatible wrapper returning (bool, int)
func (p *ReportParser) ParseCompat(message string) (bool, int) {
	result := p.Parse(message)
	return result.IsReport, result.Pages
}

func (p *ReportParser) extractPages(message string) int {
	// Range: "halaman 2-5", "5-10 hal", "hal 10 sampai 15"
	rangePattern := `(?i)(?:halaman|hal|hlm)\s*(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)`
	rangePatternRev := `(?i)(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)\s*(?:halaman|hal|hlm)`
	
	reRange := regexp.MustCompile(rangePattern)
	if m := reRange.FindStringSubmatch(message); len(m) > 2 {
		start, _ := strconv.Atoi(m[1])
		end, _ := strconv.Atoi(m[2])
		if end >= start {
			return end - start + 1
		}
	}

	reRangeRev := regexp.MustCompile(rangePatternRev)
	if m := reRangeRev.FindStringSubmatch(message); len(m) > 2 {
		start, _ := strconv.Atoi(m[1])
		end, _ := strconv.Atoi(m[2])
		if end >= start {
			return end - start + 1
		}
	}

	// Single: "5 halaman", "hal 10"
	patterns := []string{
		`(\d+)\s*(?:halaman|hal|hlm)\b`,
		`\b(?:halaman|hal|hlm)\s*(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			if val, err := strconv.Atoi(matches[1]); err == nil {
				return val
			}
		}
	}
	return 0
}

func (p *ReportParser) extractJuz(message string) int {
	// Range: "juz 1-2", "juz 1 s/d 3"
	rangePattern := `(?i)juz\s*(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)`
	reRange := regexp.MustCompile(rangePattern)
	if m := reRange.FindStringSubmatch(message); len(m) > 2 {
		start, _ := strconv.Atoi(m[1])
		end, _ := strconv.Atoi(m[2])
		if end >= start {
			return (end - start + 1) * 20
		}
	}

	// Single: "1 juz", "juz 30"
	patterns := []string{
		`(\d+)\s*juz\b`,
		`\bjuz\s*(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			if val, err := strconv.Atoi(matches[1]); err == nil {
				return val * 20
			}
		}
	}

	if strings.Contains(message, "juz") {
		return 20
	}

	return 0
}

func (p *ReportParser) extractSurahAyah(message string) ParseResult {
	// Separator pattern for ranges
	sep := `\s*(?:-+|s/d|sampai|sd|ke|dari)\s*`

	// Multiple patterns to try, ordered from most specific to least specific
	patterns := []string{
		// Pattern 1: "surat/surah <name> [ayat] <start>[-<end>]"
		`(?i)(?:surat|surah)\s+([a-z\s'-]+?)(?:\s+(?:ayat\s+)?(\d+)(?:` + sep + `(\d+))?)?(?:\s|$)`,
		// Pattern 2: "<name> ayat <start>[-<end>]" (without surat prefix, but requires 'ayat')
		`(?i)\b([a-z][a-z\s'-]{2,}?)\s+ayat\s+(\d+)(?:` + sep + `(\d+))?(?:\s|$)`,
		// Pattern 3: "baca/tilawah <name> <start>[-<end>]"
		`(?i)(?:baca|tilawah|membaca)\s+([a-z][a-z\s'-]{2,}?)\s+(\d+)(?:` + sep + `(\d+))?(?:\s|$)`,
		// Pattern 4: Loose match for common surahs followed directly by numbers (e.g., "Yasin 1-10", "Al-Baqarah 255")
		// Requires at least a range or a number > 1 to avoid matching everything
		`(?i)\b([a-z][a-z\s'-]{2,}?)\s+(\d+)(?:` + sep + `(\d+))(?:\s|$)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(message, -1)

		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			fullCapturedName := strings.TrimSpace(match[1])
			
			// Try to find surah number, potentially stripping leading words
			// e.g., "alhamdulillah yasiin" -> "yasiin" -> 36
			surahNum := 0
			words := strings.Fields(fullCapturedName)
			for i := 0; i < len(words); i++ {
				candidate := strings.Join(words[i:], " ")
				surahNum = FindSurahNumber(candidate)
				if surahNum > 0 {
					break
				}
			}

			if surahNum > 0 {
				officialName := GetSurahName(surahNum)
				startAyahStr := ""
				endAyahStr := ""
				if len(match) > 2 {
					startAyahStr = match[2]
				}
				if len(match) > 3 {
					endAyahStr = match[3]
				}

				var startAyah, endAyah int

				if startAyahStr != "" {
					startAyah, _ = strconv.Atoi(startAyahStr)
					if endAyahStr != "" {
						endAyah, _ = strconv.Atoi(endAyahStr)
					} else {
						endAyah = startAyah
					}
				} else {
					// No ayahs mentioned, get full surah pages
					startAyah = 1
					endAyah = getSurahMaxAyahs(surahNum)
				}

				// Calculate page difference
				startPage := getAyahPage(surahNum, startAyah)
				endPage := getAyahPage(surahNum, endAyah)

				pages := 1
				if startPage > 0 && endPage > 0 {
					pages = endPage - startPage + 1
					if pages < 1 {
						pages = 1
					}
				}

				return ParseResult{
					IsReport:   true,
					Pages:      pages,
					SurahName:  officialName,
					StartAyah:  startAyah,
					EndAyah:    endAyah,
					ReportType: "surah",
				}
			}
		}
	}

	return ParseResult{}
}

// FormatSurahInfo returns a formatted string like "Surat Al-Baqarah ayat 1-30 (5 halaman)"
func (r *ParseResult) FormatSurahInfo() string {
	if r.SurahName == "" {
		return ""
	}
	if r.StartAyah > 0 && r.EndAyah > 0 && r.StartAyah != r.EndAyah {
		return fmt.Sprintf("Surat %s ayat %d - %d", r.SurahName, r.StartAyah, r.EndAyah)
	}
	if r.StartAyah > 0 {
		return fmt.Sprintf("Surat %s ayat %d", r.SurahName, r.StartAyah)
	}
	return fmt.Sprintf("Surat %s (seluruh surat)", r.SurahName)
}

// getAyahPage looks up the Mushaf page for a given surah:ayah from the local dataset
func getAyahPage(surahNum, ayahNum int) int {
	if surahMap, ok := MushafPages[surahNum]; ok {
		if page, ok := surahMap[ayahNum]; ok {
			return page
		}
	}
	return 0
}

// getSurahMaxAyahs returns the total number of ayahs in a surah from the local dataset
func getSurahMaxAyahs(surahNum int) int {
	if surahMap, ok := MushafPages[surahNum]; ok {
		maxAyah := 0
		for ayah := range surahMap {
			if ayah > maxAyah {
				maxAyah = ayah
			}
		}
		return maxAyah
	}
	return 1
}
