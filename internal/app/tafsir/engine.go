package tafsir

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// maxSummaryRunes is the maximum length of the summarized tafsir text.
const maxSummaryRunes = 800

// AyahTafsir holds the tafsir text for a single verse.
type AyahTafsir struct {
	Ayah int    `json:"ayat"`
	Teks string `json:"teks"`
}

// Engine fetches and caches tafsir (Tafsir Kemenag via equran.id)
// and produces short summaries that can be shown in report replies.
type Engine struct {
	client *http.Client
	mu     sync.Mutex
	cache  map[int][]AyahTafsir
}

func NewEngine() *Engine {
	return &Engine{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
		cache: make(map[int][]AyahTafsir),
	}
}

// Summarize returns a short Indonesian summary of the tafsir for a given
// surah:ayah. It returns an error when the data is unavailable so callers
// can gracefully skip the summary.
func (e *Engine) Summarize(ctx context.Context, surahNum, ayahNum int) (string, error) {
	ayats, err := e.fetchSurahTafsir(ctx, surahNum)
	if err != nil {
		return "", err
	}

	var entry *AyahTafsir
	for i := range ayats {
		if ayats[i].Ayah == ayahNum {
			entry = &ayats[i]
			break
		}
	}
	if entry == nil {
		return "", fmt.Errorf("tafsir for verse %d:%d not found", surahNum, ayahNum)
	}

	return summarize(entry.Teks), nil
}

// SummarizeRange returns a short Indonesian summary covering the tafsir of an
// ayah range within a surah (e.g. a whole surah reading). For a single ayah it
// behaves like Summarize.
func (e *Engine) SummarizeRange(ctx context.Context, surahNum, startAyah, endAyah int) (string, error) {
	if endAyah < startAyah {
		startAyah, endAyah = endAyah, startAyah
	}
	if startAyah == endAyah {
		return e.Summarize(ctx, surahNum, startAyah)
	}

	ayats, err := e.fetchSurahTafsir(ctx, surahNum)
	if err != nil {
		return "", err
	}

	var sentences []string
	for _, a := range ayats {
		if a.Ayah < startAyah || a.Ayah > endAyah {
			continue
		}
		clean := htmlTagRegex.ReplaceAllString(a.Teks, "")
		clean = strings.Join(strings.Fields(clean), " ")
		if clean == "" {
			continue
		}

		parts := splitSentences(clean)
		if a.Ayah == startAyah {
			// Keep a fuller intro for the first verse of the range.
			if len(parts) > 2 {
				parts = parts[:2]
			}
			sentences = append(sentences, parts...)
			continue
		}

		// For the rest, take only the first meaningful sentence, skipping
		// cross-reference notes like "Lihat Tafsir Alif Lam Mim...".
		for _, s := range parts {
			if lihatTafsirRegex.MatchString(s) {
				continue
			}
			sentences = append(sentences, s)
			break
		}
	}

	if len(sentences) == 0 {
		return "", fmt.Errorf("tafsir for range %d:%d-%d not found", surahNum, startAyah, endAyah)
	}

	return truncateAtWordBoundary(strings.Join(sentences, " "), maxSummaryRunes), nil
}

// fetchSurahTafsir fetches and caches the full tafsir of a surah.
func (e *Engine) fetchSurahTafsir(ctx context.Context, surahNum int) ([]AyahTafsir, error) {
	e.mu.Lock()
	cached, ok := e.cache[surahNum]
	e.mu.Unlock()
	if ok {
		return cached, nil
	}

	url := fmt.Sprintf("https://equran.id/api/v2/tafsir/%d", surahNum)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tafsir API error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data struct {
			Tafsir []AyahTafsir `json:"tafsir"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Data.Tafsir) == 0 {
		return nil, fmt.Errorf("no tafsir data for surah %d", surahNum)
	}

	e.mu.Lock()
	e.cache[surahNum] = result.Data.Tafsir
	e.mu.Unlock()

	return result.Data.Tafsir, nil
}

var htmlTagRegex = regexp.MustCompile("<[^>]*>")

var lihatTafsirRegex = regexp.MustCompile(`(?i)^lihat\s+tafsir`)

// summarize turns the full tafsir text into a short summary (up to a few
// sentences so the report reply stays readable).
func summarize(teks string) string {
	clean := htmlTagRegex.ReplaceAllString(teks, "")
	clean = strings.Join(strings.Fields(clean), " ")
	if clean == "" {
		return ""
	}

	sentences := splitSentences(clean)
	if len(sentences) > 4 {
		sentences = sentences[:4]
	}

	summary := strings.Join(sentences, " ")
	return truncateAtWordBoundary(summary, maxSummaryRunes)
}

// splitSentences splits text into individual printable sentences.
func splitSentences(text string) []string {
	var sentences []string
	start := 0
	for i := 0; i < len(text); i++ {
		if text[i] != '.' && text[i] != '!' && text[i] != '?' {
			continue
		}
		if i+1 < len(text) {
			next := text[i+1]
			if next != '.' && next != ' ' && next != ')' && next != '"' {
				continue
			}
		}
		if s := strings.TrimSpace(text[start : i+1]); s != "" {
			sentences = append(sentences, s)
		}
		start = i + 1
	}
	if s := strings.TrimSpace(text[start:]); s != "" {
		sentences = append(sentences, s)
	}
	return sentences
}

// truncateAtWordBoundary shortens s to at most maxChars runes, breaking at a
// word boundary and appending an ellipsis when truncated.
func truncateAtWordBoundary(s string, maxChars int) string {
	runes := []rune(s)
	if len(runes) <= maxChars {
		return strings.TrimSpace(s)
	}

	trimmed := string(runes[:maxChars])
	if idx := strings.LastIndex(trimmed, " "); idx > maxChars/2 {
		trimmed = trimmed[:idx]
	}
	return strings.TrimSpace(trimmed) + "…"
}