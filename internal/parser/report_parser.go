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

// Parse determines if the message contains a valid report and returns a slice of ParseResult
func (p *ReportParser) Parse(message string) []ParseResult {
	// Must contain alhamdulillah (flexible match)
	alhamdulillahRegex := regexp.MustCompile(`(?i)#?al[ -]?hamdu?[ -]?l+il+a+h`)
	if !alhamdulillahRegex.MatchString(message) {
		return nil
	}

	var results []ParseResult
	
	// We use a copy of the message that we'll modify (consume)
	workMsg := strings.ToLower(message)

	// 1. Try surah range (e.g. Al-Ashr - An-Nas)
	rangeResults := p.extractSurahRange(&workMsg)
	if len(rangeResults) > 0 {
		results = append(results, rangeResults...)
	}

	// 2. Try surah + ayah (individual surahs)
	surahResults := p.extractSurahAyah(&workMsg)
	if len(surahResults) > 0 {
		results = append(results, surahResults...)
	}

	// 2. Try juz (can have multiple)
	juzResults := p.extractJuz(&workMsg)
	if len(juzResults) > 0 {
		results = append(results, juzResults...)
	}

	// 3. Try halaman
	// Always try extractPages even if we found surah/juz, but it will use the remaining workMsg
	if pages := p.extractPages(&workMsg); pages > 0 {
		results = append(results, ParseResult{IsReport: true, Pages: pages, ReportType: "halaman"})
	}

	// 4. Default: assume 1 page if alhamdulillah detected but no specific info found
	if len(results) == 0 {
		results = append(results, ParseResult{IsReport: true, Pages: 1, ReportType: "default"})
	}

	return results
}

// ParseCompat is a backward-compatible wrapper returning (bool, int)
// If multiple results, it returns the total pages.
func (p *ReportParser) ParseCompat(message string) (bool, int) {
	results := p.Parse(message)
	if len(results) == 0 {
		return false, 0
	}
	totalPages := 0
	for _, r := range results {
		totalPages += r.Pages
	}
	return true, totalPages
}

func (p *ReportParser) extractSurahRange(message *string) []ParseResult {
	// Pattern for range separators. Require at least one space around hyphen.
	sepPattern := `\s+(-+|sampai|sd|s/d|ke|dari)\s+`
	re := regexp.MustCompile(sepPattern)
	
	var results []ParseResult
	
	// Use FindAllStringIndex to avoid infinite loop when we skip masking
	matches := re.FindAllStringIndex(*message, -1)
	
	// Process matches in reverse to avoid offset issues if we mask
	for i := len(matches) - 1; i >= 0; i-- {
		loc := matches[i]
		beforePart := (*message)[:loc[0]]
		afterPart := (*message)[loc[1]:]

		wordsBefore := strings.Fields(beforePart)
		wordsAfter := strings.Fields(afterPart)

		s1, s2 := 0, 0
		len1, len2 := 0, 0 

		for i := 1; i <= 3 && i <= len(wordsBefore); i++ {
			candidate := strings.Join(wordsBefore[len(wordsBefore)-i:], " ")
			if num := FindSurahNumber(candidate); num > 0 {
				s1 = num
				len1 = i
			}
		}

		for i := 1; i <= 3 && i <= len(wordsAfter); i++ {
			candidate := strings.Join(wordsAfter[:i], " ")
			if num := FindSurahNumber(candidate); num > 0 {
				s2 = num
				len2 = i
			}
		}

		if s1 <= 0 || s2 <= 0 || s2 < s1 {
			continue
		}

		// Check for numeric range false positives
		s1Text_candidate := strings.Join(wordsBefore[len(wordsBefore)-len1:], " ")
		s2Text_candidate := strings.Join(wordsAfter[:len2], " ")
		isNumeric1 := regexp.MustCompile(`^\d+$`).MatchString(s1Text_candidate)
		isNumeric2 := regexp.MustCompile(`^\d+$`).MatchString(s2Text_candidate)
		// Check for optional "ayat <number>" after s2Text
		endAyah := 0
		lenAyah := 0
		if len(wordsAfter) > len2 {
			// Check if next word is "ayat"
			if strings.ToLower(wordsAfter[len2]) == "ayat" {
				if len(wordsAfter) > len2+1 {
					if num, err := strconv.Atoi(wordsAfter[len2+1]); err == nil {
						endAyah = num
						lenAyah = 2
					}
				}
			} else {
				// Maybe just a number? e.g. "Al-Baqarah 100" (but wait, we need to be careful with "Al-Baqarah 100-110")
				// Usually in ranges like "Fatihah sampai Baqarah 100", people might omit 'ayat'.
				// But let's require 'ayat' for now to be safe, or check if it's a number followed by nothing else.
				if num, err := strconv.Atoi(wordsAfter[len2]); err == nil {
					endAyah = num
					lenAyah = 1
				}
			}
		}

		if isNumeric1 && isNumeric2 {
			// ... (numeric range check same as before)
			hasPrefix := false
			if len(wordsBefore) > len1 {
				lastWord := strings.ToLower(wordsBefore[len(wordsBefore)-len1-1])
				if lastWord == "surat" || lastWord == "surah" {
					hasPrefix = true
				}
			}
			if !hasPrefix {
				continue
			}
		}

		// ... (ayah range false positive check same as before)
		beforeTrim := strings.TrimSpace(beforePart)
		if len(beforeTrim) > 0 {
			lastChar := beforeTrim[len(beforeTrim)-1]
			if lastChar >= '0' && lastChar <= '9' {
				continue
			}
		}

		startPage := getAyahPage(s1, 1)
		
		finalEndAyah := getSurahMaxAyahs(s2)
		if endAyah > 0 && endAyah <= finalEndAyah {
			finalEndAyah = endAyah
		}
		endPage := getAyahPage(s2, finalEndAyah)
		
		pages := 1
		if startPage > 0 && endPage > 0 {
			pages = endPage - startPage + 1
			if pages < 1 { pages = 1 }
		}

		surahName := fmt.Sprintf("%s - %s", GetSurahName(s1), GetSurahName(s2))
		if endAyah > 0 {
			surahName = fmt.Sprintf("%s - %s ayat %d", GetSurahName(s1), GetSurahName(s2), endAyah)
		}

		results = append(results, ParseResult{
			IsReport:   true,
			Pages:      pages,
			SurahName:  surahName,
			ReportType: "surah",
		})

		// Mask the match (s1 text to s2 text + optional ayah index)
		idx1 := strings.LastIndex(beforePart, s1Text_candidate)
		
		// Recalculate end index including ayah if found
		endWords := wordsAfter[:len2+lenAyah]
		// Note: Joining might lose original spacing but strings.Index in afterPart should be fine
		// if we use a more robust way to find the end position.
		
		// Find end of the wordsAfter section we want to mask
		// A simple way is to find the index of the last word's end.
		lastWordInMatch := endWords[len(endWords)-1]
		// Start searching for last word after we've seen all preceding words
		searchOffset := 0
		for _, w := range endWords[:len(endWords)-1] {
			if pos := strings.Index(afterPart[searchOffset:], w); pos != -1 {
				searchOffset += pos + len(w)
			}
		}
		idxOfLastWord := strings.Index(afterPart[searchOffset:], lastWordInMatch)
		
		if idx1 != -1 && idxOfLastWord != -1 {
			totalMatchStart := idx1
			totalMatchEnd := loc[1] + searchOffset + idxOfLastWord + len(lastWordInMatch)
			*message = (*message)[:totalMatchStart] + strings.Repeat(" ", totalMatchEnd-totalMatchStart) + (*message)[totalMatchEnd:]
		} else {
			*message = (*message)[:loc[0]] + strings.Repeat(" ", loc[1]-loc[0]) + (*message)[loc[1]:]
		}
	}
	return results
}

func (p *ReportParser) extractPages(message *string) int {
	total := 0
	patterns := []string{
		`(?i)(?:halaman|hal|hlm)\s*(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)`,
		`(?i)(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)\s*(?:halaman|hal|hlm)`,
		`(\d+)/(\d+)\s*(?:halaman|hal|hlm)\b`,
		`(\d+(?:\.\d+)?)\s*(?:halaman|hal|hlm)\b`,
		`\b(?:halaman|hal|hlm)\s*(\d+(?:\.\d+)?)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(*message, -1)
		for _, loc := range matches {
			match := (*message)[loc[0]:loc[1]]
			// Check if already masked
			if strings.TrimSpace(match) == "" {
				continue
			}

			if strings.Contains(pattern, `(\d+)/(\d+)`) {
				num, _ := strconv.ParseFloat((*message)[loc[2]:loc[3]], 64)
				den, _ := strconv.ParseFloat((*message)[loc[4]:loc[5]], 64)
				if den > 0 {
					total += int(num / den)
				}
			} else if strings.Contains(pattern, `(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)`) {
				start, _ := strconv.Atoi((*message)[loc[2]:loc[3]])
				end, _ := strconv.Atoi((*message)[loc[4]:loc[5]])
				if end >= start {
					total += end - start + 1
				}
			} else {
				if val, err := strconv.ParseFloat((*message)[loc[2]:loc[3]], 64); err == nil {
					total += int(val)
				}
			}
			*message = (*message)[:loc[0]] + strings.Repeat(" ", loc[1]-loc[0]) + (*message)[loc[1]:]
		}
	}
	return total
}

func (p *ReportParser) extractJuz(message *string) []ParseResult {
	var results []ParseResult
	patterns := []string{
		`(?i)juz\s*(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)`,
		`(\d+)/(\d+)\s*juz\b`,
		`(\d+(?:\.\d+)?)\s*juz\b`,
		`(?i)\bjuz\s*(\d+)`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(*message, -1)
		for _, loc := range matches {
			match := (*message)[loc[0]:loc[1]]
			if strings.TrimSpace(match) == "" {
				continue
			}

			if strings.Contains(pattern, `(\d+)/(\d+)`) {
				num, _ := strconv.ParseFloat((*message)[loc[2]:loc[3]], 64)
				den, _ := strconv.ParseFloat((*message)[loc[4]:loc[5]], 64)
				if den > 0 {
					results = append(results, ParseResult{IsReport: true, Pages: int(num / den * 20), ReportType: "juz"})
				}
			} else if strings.Contains(pattern, `(\d+)\s*(?:-+|s/d|sampai|sd|ke|dari)\s*(\d+)`) {
				start, _ := strconv.Atoi((*message)[loc[2]:loc[3]])
				end, _ := strconv.Atoi((*message)[loc[4]:loc[5]])
				if end >= start {
					results = append(results, ParseResult{IsReport: true, Pages: (end - start + 1) * 20, ReportType: "juz"})
				}
			} else if strings.Contains(pattern, `(\d+(?:\.\d+)?)\s*juz\b`) {
				if val, err := strconv.ParseFloat((*message)[loc[2]:loc[3]], 64); err == nil {
					results = append(results, ParseResult{IsReport: true, Pages: int(val * 20), ReportType: "juz"})
				}
			} else {
				results = append(results, ParseResult{IsReport: true, Pages: 20, ReportType: "juz"})
			}
			*message = (*message)[:loc[0]] + strings.Repeat(" ", loc[1]-loc[0]) + (*message)[loc[1]:]
		}
	}
	return results
}

func (p *ReportParser) extractSurahAyah(message *string) []ParseResult {
	sep := `\s*(?:-+|s/d|sampai|sd|ke|dari)\s*`
	patterns := []string{
		`(?i)(?:surat|surah)\s+([a-z\s'-]+?)(?:\s+(?:ayat\s+)?(\d+)(?:` + sep + `(\d+))?)?(?:\b|$)`,
		`(?i)\b([a-z][a-z\s'-]{2,}?)\s+ayat\s+(\d+)(?:` + sep + `(\d+))?(?:\b|$)`,
		`(?i)(?:baca|tilawah|membaca)\s+([a-z][a-z\s'-]{2,}?)\s+(\d+)(?:` + sep + `(\d+))?(?:\b|$)`,
		`(?i)\b([a-z][a-z\s'-]{2,}?)\s+(\d+)(?:` + sep + `(\d+))(?:\b|$)`,
		`(?i)\b([a-z][a-z\s'-]{2,}?)\b`,
	}

	var results []ParseResult
	seenSurahs := make(map[string]bool)

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatchIndex(*message, -1)
		for _, loc := range matches {
			if strings.TrimSpace((*message)[loc[0]:loc[1]]) == "" {
				continue
			}

			fullCapturedName := strings.TrimSpace((*message)[loc[2]:loc[3]])
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
				if seenSurahs[officialName] {
					*message = (*message)[:loc[0]] + strings.Repeat(" ", loc[1]-loc[0]) + (*message)[loc[1]:]
					continue
				}

				startAyahStr := ""
				endAyahStr := ""
				if len(loc) > 4 && loc[4] != -1 { startAyahStr = (*message)[loc[4]:loc[5]] }
				if len(loc) > 6 && loc[6] != -1 { endAyahStr = (*message)[loc[6]:loc[7]] }

				var startAyah, endAyah int
				maxAyah := getSurahMaxAyahs(surahNum)
				if startAyahStr != "" {
					startAyah, _ = strconv.Atoi(startAyahStr)
					if endAyahStr != "" { endAyah, _ = strconv.Atoi(endAyahStr) } else { endAyah = startAyah }
				} else { startAyah = 1; endAyah = maxAyah }

				if startAyah < 1 { startAyah = 1 }
				if startAyah > maxAyah { startAyah = maxAyah }
				if endAyah < startAyah { endAyah = startAyah }
				if endAyah > maxAyah { endAyah = maxAyah }

				startPage := getAyahPage(surahNum, startAyah)
				endPage := getAyahPage(surahNum, endAyah)
				pages := 1
				if startPage > 0 && endPage > 0 {
					pages = endPage - startPage + 1
					if pages < 1 { pages = 1 }
				}

				results = append(results, ParseResult{
					IsReport: true, Pages: pages, SurahName: officialName,
					StartAyah: startAyah, EndAyah: endAyah, ReportType: "surah",
				})
				seenSurahs[officialName] = true
				*message = (*message)[:loc[0]] + strings.Repeat(" ", loc[1]-loc[0]) + (*message)[loc[1]:]
			}
		}
	}
	return results
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
		// fallback: if ayahNum is out of range, try to get closest valid ayah
		maxAyah := getSurahMaxAyahs(surahNum)
		if ayahNum <= 0 {
			return surahMap[1]
		}
		if ayahNum > maxAyah {
			return surahMap[maxAyah]
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
