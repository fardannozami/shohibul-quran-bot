package parser

import (
	"strings"
)

// SurahList contains the standard Arabic-transliterated names of all 114 surahs
var SurahList = []string{
	"Al-Fatihah", "Al-Baqarah", "Ali 'Imran", "An-Nisa'", "Al-Ma'idah", "Al-An'am", "Al-A'raf", "Al-Anfal",
	"At-Taubah", "Yunus", "Hud", "Yusuf", "Ar-Ra'd", "Ibrahim", "Al-Hijr", "An-Nahl", "Al-Isra'", "Al-Kahfi",
	"Maryam", "Taha", "Al-Anbiya'", "Al-Hajj", "Al-Mu'minun", "An-Nur", "Al-Furqan", "Asy-Syu'ara'", "An-Naml",
	"Al-Qasas", "Al-'Ankabut", "Ar-Rum", "Luqman", "As-Sajdah", "Al-Ahzab", "Saba'", "Fatir", "Yasin",
	"As-Saffat", "Sad", "Az-Zumar", "Ghafir", "Fussilat", "Asy-Syura", "Az-Zukhruf", "Ad-Dukhan", "Al-Jasiyah",
	"Al-Ahqaf", "Muhammad", "Al-Fath", "Al-Hujurat", "Qaf", "Az-Zariyat", "At-Tur", "An-Najm", "Al-Qamar",
	"Ar-Rahman", "Al-Waqi'ah", "Al-Hadid", "Al-Mujadalah", "Al-Hasyr", "Al-Mumtahanah", "As-Saff", "Al-Jumu'ah",
	"Al-Munafiqun", "At-Tagabun", "At-Talaq", "At-Tahrim", "Al-Mulk", "Al-Qalam", "Al-Haqqah", "Al-Ma'arij",
	"Nuh", "Al-Jinn", "Al-Muzzammil", "Al-Muddassir", "Al-Qiyamah", "Al-Insan", "Al-Mursalat", "An-Naba'",
	"An-Nazi'at", "'Abasa", "At-Takwir", "Al-Infitar", "Al-Mutaffifin", "Al-Insyiqaq", "Al-Buruj", "At-Tariq",
	"Al-A'la", "Al-Ghasyiyah", "Al-Fajr", "Al-Balad", "Asy-Syams", "Al-Lail", "Ad-Duha", "Asy-Syarh", "At-Tin",
	"Al-'Alaq", "Al-Qadr", "Al-Bayyinah", "Az-Zalzalah", "Al-'Adiyat", "Al-Qari'ah", "At-Takasur", "Al-'Asr",
	"Al-Humazah", "Al-Fil", "Quraisy", "Al-Ma'un", "Al-Kausar", "Al-Kafirun", "An-Nasr", "Al-Lahab", "Al-Ikhlas",
	"Al-Falaq", "An-Nas",
}

// surahAliases maps informal/Indonesian spellings to surah number (1-indexed)
// This handles common typos and informal names used in WhatsApp groups
var surahAliases = map[string]int{
	// Surah 1
	"fatihah": 1, "alfatihah": 1, "fatiha": 1, "alfatiha": 1, "fathia": 1,
	// Surah 2
	"baqarah": 2, "albaqarah": 2, "baqoroh": 2, "albaqoroh": 2, "baqara": 2, "bakara": 2, "bakarah": 2, "albaqoro": 2,
	// Surah 3
	"imran": 3, "aliimran": 3, "imron": 3, "aliimron": 3,
	// Surah 4
	"nisa": 4, "annisa": 4, "nisaa": 4, "annisaa": 4,
	// Surah 5
	"maidah": 5, "almaidah": 5, "maaidah": 5, "maidoh": 5,
	// Surah 6
	"anam": 6, "alanam": 6, "anaam": 6,
	// Surah 7
	"araf": 7, "alaraf": 7, "a'raf": 7,
	// Surah 8
	"anfal": 8, "alanfal": 8,
	// Surah 9
	"taubah": 9, "attaubah": 9, "taubat": 9, "tawbah": 9,
	// Surah 10
	"yunus": 10,
	// Surah 11
	"hud": 11, "huud": 11,
	// Surah 12
	"yusuf": 12, "yusup": 12, "jusuf": 12,
	// Surah 13
	"rad": 13, "arrad": 13, "raad": 13,
	// Surah 14
	"ibrahim": 14,
	// Surah 15
	"hijr": 15, "alhijr": 15,
	// Surah 16
	"nahl": 16, "annahl": 16,
	// Surah 17
	"isra": 17, "alisra": 17, "isro": 17,
	// Surah 18
	"kahfi": 18, "alkahfi": 18, "kahf": 18, "alkahf": 18,
	// Surah 19
	"maryam": 19, "mariam": 19,
	// Surah 20
	"taha": 20, "toha": 20, "thaha": 20,
	// Surah 21
	"anbiya": 21, "alanbiya": 21,
	// Surah 22
	"hajj": 22, "alhajj": 22, "haji": 22,
	// Surah 23
	"muminun": 23, "almuminun": 23, "mukminun": 23,
	// Surah 24
	"nur": 24, "annur": 24, "nuur": 24,
	// Surah 25
	"furqan": 25, "alfurqan": 25, "furqon": 25,
	// Surah 26
	"syuara": 26, "asysyuara": 26,
	// Surah 27
	"naml": 27, "annaml": 27,
	// Surah 28
	"qasas": 28, "alqasas": 28, "qosos": 28,
	// Surah 29
	"ankabut": 29, "alankabut": 29,
	// Surah 30
	"rum": 30, "arrum": 30,
	// Surah 31
	"luqman": 31, "lukman": 31,
	// Surah 32
	"sajdah": 32, "assajdah": 32, "sajadah": 32,
	// Surah 33
	"ahzab": 33, "alahzab": 33,
	// Surah 34
	"saba": 34,
	// Surah 35
	"fatir": 35, "faatir": 35,
	// Surah 36
	"yasin": 36, "yasiin": 36, "yaasin": 36, "yaasiin": 36, "ysin": 36,
	// Surah 37
	"saffat": 37, "assaffat": 37, "shaffat": 37,
	// Surah 38
	"sad": 38, "shad": 38, "shaad": 38,
	// Surah 39
	"zumar": 39, "azzumar": 39,
	// Surah 40
	"ghafir": 40, "gafir": 40, "ghofir": 40,
	// Surah 41
	"fussilat": 41, "fushshilat": 41,
	// Surah 42
	"syura": 42, "asysyura": 42, "syuura": 42, "syuro": 42, "asysyuro": 42, "suaro": 42,
	// Surah 43
	"zukhruf": 43, "azzukhruf": 43,
	// Surah 44
	"dukhan": 44, "addukhan": 44, "dukhon": 44,
	// Surah 45
	"jasiyah": 45, "aljasiyah": 45, "jatsiyah": 45,
	// Surah 46
	"ahqaf": 46, "alahqaf": 46,
	// Surah 47
	"muhammad": 47,
	// Surah 48
	"fath": 48, "alfath": 48, "alfathu": 48,
	// Surah 49
	"hujurat": 49, "alhujurat": 49, "hujurot": 49,
	// Surah 50
	"qaf": 50, "qaaf": 50,
	// Surah 51
	"zariyat": 51, "azzariyat": 51, "dzariyat": 51,
	// Surah 52
	"tur": 52, "attur": 52, "thur": 52,
	// Surah 53
	"najm": 53, "annajm": 53,
	// Surah 54
	"qamar": 54, "alqamar": 54, "qomar": 54,
	// Surah 55
	"rahman": 55, "arrahman": 55, "rohman": 55,
	// Surah 56
	"waqiah": 56, "alwaqiah": 56, "waqi'ah": 56, "waqia": 56, "wakiah": 56,
	// Surah 57
	"hadid": 57, "alhadid": 57,
	// Surah 58
	"mujadalah": 58, "almujadalah": 58, "mujadilah": 58,
	// Surah 59
	"hasyr": 59, "alhasyr": 59, "hasyir": 59,
	// Surah 60
	"mumtahanah": 60, "almumtahanah": 60,
	// Surah 61
	"saff": 61, "assaff": 61, "shaf": 61, "shoff": 61,
	// Surah 62
	"jumuah": 62, "aljumuah": 62, "jumah": 62,
	// Surah 63
	"munafiqun": 63, "almunafiqun": 63, "munafiqin": 63,
	// Surah 64
	"tagabun": 64, "attagabun": 64, "tagobun": 64,
	// Surah 65
	"talaq": 65, "attalaq": 65, "tholaq": 65,
	// Surah 66
	"tahrim": 66, "attahrim": 66,
	// Surah 67
	"mulk": 67, "almulk": 67,
	// Surah 68
	"qalam": 68, "alqalam": 68,
	// Surah 69
	"haqqah": 69, "alhaqqah": 69,
	// Surah 70
	"maarij": 70, "almaarij": 70,
	// Surah 71
	"nuh": 71, "nuuh": 71, "nooh": 71,
	// Surah 72
	"jinn": 72, "aljinn": 72, "jin": 72,
	// Surah 73
	"muzzammil": 73, "almuzzammil": 73, "muzammil": 73,
	// Surah 74
	"muddassir": 74, "almuddassir": 74, "mudatsir": 74,
	// Surah 75
	"qiyamah": 75, "alqiyamah": 75, "qiyamahi": 75,
	// Surah 76
	"insan": 76, "alinsan": 76,
	// Surah 77
	"mursalat": 77, "almursalat": 77,
	// Surah 78
	"naba": 78, "annaba": 78,
	// Surah 79
	"naziat": 79, "annaziat": 79,
	// Surah 80
	"abasa": 80,
	// Surah 81
	"takwir": 81, "attakwir": 81,
	// Surah 82
	"infitar": 82, "alinfitar": 82,
	// Surah 83
	"mutaffifin": 83, "almutaffifin": 83, "tatfif": 83,
	// Surah 84
	"insyiqaq": 84, "alinsyiqaq": 84, "insiqaq": 84,
	// Surah 85
	"buruj": 85, "alburuj": 85,
	// Surah 86
	"tariq": 86, "attariq": 86,
	// Surah 87
	"ala": 87, "alala": 87,
	// Surah 88
	"ghasyiyah": 88, "alghasyiyah": 88, "gosyiyah": 88,
	// Surah 89
	"fajr": 89, "alfajr": 89, "fajar": 89,
	// Surah 90
	"balad": 90, "albalad": 90,
	// Surah 91
	"syams": 91, "asysyams": 91,
	// Surah 92
	"lail": 92, "allail": 92,
	// Surah 93
	"duha": 93, "adduha": 93, "dhuha": 93,
	// Surah 94
	"syarh": 94, "asysyarh": 94, "insyirah": 94, "alinsyirah": 94,
	// Surah 95
	"tin": 95, "attin": 95,
	// Surah 96
	"alaq": 96, "alalaq": 96,
	// Surah 97
	"qadr": 97, "alqadr": 97, "qodr": 97,
	// Surah 98
	"bayyinah": 98, "albayyinah": 98,
	// Surah 99
	"zalzalah": 99, "azzalzalah": 99,
	// Surah 100
	"adiyat": 100, "aladiyat": 100,
	// Surah 101
	"qariah": 101, "alqariah": 101,
	// Surah 102
	"takasur": 102, "attakasur": 102, "takatsur": 102,
	// Surah 103
	"asr": 103, "alasr": 103, "ashr": 103, "alashr": 103,
	// Surah 104
	"humazah": 104, "alhumazah": 104,
	// Surah 105
	"fil": 105, "alfil": 105,
	// Surah 106
	"quraisy": 106, "quraish": 106,
	// Surah 107
	"maun": 107, "almaun": 107,
	// Surah 108
	"kausar": 108, "alkausar": 108, "kautsar": 108, "kautsari": 108, "kauthar": 108,
	// Surah 109
	"kafirun": 109, "alkafirun": 109, "kafirin": 109,
	// Surah 110
	"nasr": 110, "annasr": 110, "nashr": 110,
	// Surah 111
	"lahab": 111, "allahab": 111, "masad": 111, "almasad": 111,
	// Surah 112
	"ikhlas": 112, "alikhlas": 112, "iklas": 112, "aliklas": 112,
	// Surah 113
	"falaq": 113, "alfalaq": 113,
	// Surah 114
	"nas": 114, "annas": 114,
}

func sanitizeSurahName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, "-", "")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, "'", "")
	name = strings.ReplaceAll(name, " ", "")
	return name
}

// GetSurahName returns the proper Arabic-transliterated name of a surah by number
func GetSurahName(surahNum int) string {
	if surahNum >= 1 && surahNum <= len(SurahList) {
		return SurahList[surahNum-1]
	}
	return ""
}

// FindSurahNumber tries to match a given string to a Surah number (1-114). Returns 0 if not found.
func FindSurahNumber(input string) int {
	sanitized := sanitizeSurahName(input)

	// Blacklist of common words that should never match a surah
	noiseWords := map[string]bool{
		"alhamdulillah": true, "alhamdulilah": true, "alhamdullillah": true,
		"ayat": true, "surat": true, "surah": true, "juz": true, "juzuk": true,
		"hal": true, "halaman": true, "hlm": true, "lembar": true,
		"baca": true, "tilawah": true, "membaca": true,
		"dan": true, "serta": true, "ke": true, "dari": true, "sampai": true,
		"hari": true, "ini": true, "tadi": true, "malam": true,
	}
	if noiseWords[sanitized] {
		return 0
	}

	// 1. Exact Match on Aliases (handles typos and informal names)
	if num, ok := surahAliases[sanitized]; ok {
		return num
	}

	// 2. Exact Match on Official Names
	for i, surah := range SurahList {
		if sanitizeSurahName(surah) == sanitized {
			return i + 1
		}
	}

	// 3. Prefix Match
	// We pick the shortest target that matches the prefix to be more exact
	if len(sanitized) >= 3 {
		bestPrefixMatch := 0
		shortestLen := 999

		for i, surah := range SurahList {
			official := sanitizeSurahName(surah)
			if strings.HasPrefix(official, sanitized) {
				if len(official) < shortestLen {
					shortestLen = len(official)
					bestPrefixMatch = i + 1
				} else if len(official) == shortestLen && (bestPrefixMatch == 0 || i+1 < bestPrefixMatch) {
					bestPrefixMatch = i + 1
				}
			}
		}

		for alias, num := range surahAliases {
			if strings.HasPrefix(alias, sanitized) {
				if len(alias) < shortestLen {
					shortestLen = len(alias)
					bestPrefixMatch = num
				} else if len(alias) == shortestLen && (bestPrefixMatch == 0 || num < bestPrefixMatch) {
					bestPrefixMatch = num
				}
			}
		}

		if bestPrefixMatch > 0 {
			return bestPrefixMatch
		}
	}

	// 4. Fuzzy Match (Levenshtein Distance) as last fallback
	// We only do this if the input is at least 4 characters to avoid false positives
	if len(sanitized) >= 4 {
		bestMatch := 0
		bestScore := 0.75 // High threshold for better accuracy
		bestDist := 999

		checkMatch := func(target string, num int) {
			dist := levenshteinDistance(sanitized, target)
			maxLen := len(sanitized)
			if len(target) > maxLen {
				maxLen = len(target)
			}
			score := 1.0 - (float64(dist) / float64(maxLen))

			if score > bestScore {
				bestScore = score
				bestMatch = num
				bestDist = dist
			} else if score == bestScore {
				if dist < bestDist {
					bestMatch = num
					bestDist = dist
				} else if dist == bestDist {
					// Deterministic tie-breaking: prefer lower surah number
					if bestMatch == 0 || num < bestMatch {
						bestMatch = num
					}
				}
			}
		}

		for i, surah := range SurahList {
			checkMatch(sanitizeSurahName(surah), i+1)
		}

		for alias, num := range surahAliases {
			checkMatch(alias, num)
		}

		if bestMatch > 0 {
			return bestMatch
		}
	}

	return 0
}

// similarity returns a value between 0 and 1 representing the similarity of two strings
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

// levenshteinDistance calculates the edit distance between two strings
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
