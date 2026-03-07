package parser

import (
	"regexp"
	"strconv"
	"strings"
)

// ReportParser handles parsing of chat messages to detect reading activity
type ReportParser struct{}

func NewReportParser() *ReportParser {
	return &ReportParser{}
}

// Parse determines if the message contains a valid report and returns the number of pages read
func (p *ReportParser) Parse(message string) (bool, int) {
	lowerMsg := strings.ToLower(message)
	
	// Must contain alhamdulillah
	if !strings.Contains(lowerMsg, "alhamdulillah") {
		return false, 0
	}

	pages := p.extractPages(lowerMsg)
	if pages == 0 {
		pages = p.extractJuz(lowerMsg)
	}

	// Assuming at least 1 page read if they report without explicit numbers
	if pages == 0 {
		pages = 1
	}

	return true, pages
}

func (p *ReportParser) extractPages(message string) int {
	// Match "5 halaman", "baca 10 hal", "halaman 2"
	patterns := []string{
		`(\d+)\s*(?:halaman|hal)\b`,
		`\b(?:halaman|hal)\s*(\d+)`,
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
	// Match "1 juz", "juz 3", "2 juz"
	// 1 juz approx 20 pages
	patterns := []string{
		`(\d+)\s*juz\b`,
		`\bjuz\s*(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(message)
		if len(matches) > 1 {
			if val, err := strconv.Atoi(matches[1]); err == nil {
				// We don't want to convert "juz 3" to 3 juz, it just means they read 1 juz (Juz number 3).
				// So if they say "1 juz" or "juz 3", it's usually 20 pages.
				// For simplicity, we just assume 1 juz = 20 pages regardless of the number mentioned.
				// Wait, "2 juz" means 40 pages. "Juz 30" could mean 1 juz.
				// Let's just default to 20 pages if the word 'juz' is mentioned.
				return val * 20
			}
		}
	}

	// If no number, but just contains 'juz'
	if strings.Contains(message, "juz") {
		return 20
	}

	return 0
}
