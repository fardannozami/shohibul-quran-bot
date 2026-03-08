package main

import (
	"fmt"
	"strings"
)

var SurahList = []string{
	"Al-Fatihah", "Al-Baqarah", "Ali 'Imran", "An-Nisa'", "Al-Ma'idah", "Al-An'am", "Al-A'raf", "Al-Anfal",
	"At-Taubah", "Yunus", "Hud", "Yusuf", "Ar-Ra'd", "Ibrahim", "Al-Hijr", "An-Nahl", "Al-Isra'", "Al-Kahfi",
	"Maryam", "Taha", "Al-Anbiya'", "Al-Hajj", "Al-Mu'minun", "An-Nur", "Al-Furqan", "Asy-Syu'ara'", "An-Naml",
	"Al-Qasas", "Al-'Ankabut", "Ar-Rum", "Luqman", "As-Sajdah", "Al-Ahzab", "Saba'", "Fatir", "Yasin",
}

func sanitizeSurahName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, " ", "")
	return name
}

func similarity(s1, s2 string) float64 {
	if s1 == s2 {
		return 1.0
	}
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	dist := levenshteinDistance(s1, s2)
	maxLen := len(s1)
	if len(s2) > maxLen {
		maxLen = len(s2)
	}

	return 1.0 - (float64(dist) / float64(maxLen))
}

func levenshteinDistance(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)

	row := make([]int, len2+1)
	for j := 0; j <= len2; j++ {
		row[j] = j
	}

	for i := 1; i <= len1; i++ {
		prev := i
		for j := 1; j <= len2; j++ {
			var val int
			if s1[i-1] == s2[j-1] {
				val = row[j-1]
			} else {
				min := row[j-1] // substitution
				if row[j] < min {
					min = row[j] // deletion
				}
				if prev < min {
					min = prev // insertion
				}
				val = min + 1
			}
			row[j-1] = prev
			prev = val
		}
		row[len2] = prev
	}

	return row[len2]
}

func main() {
	input := "fathia"
	sanitized := sanitizeSurahName(input)
	fmt.Printf("Input: %s, Sanitized: %s\n", input, sanitized)

	for i, surah := range SurahList {
		official := sanitizeSurahName(surah)
		score := similarity(sanitized, official)
		if score > 0.5 {
			fmt.Printf("Surah %d (%s): %f\n", i+1, official, score)
		}
	}
}
