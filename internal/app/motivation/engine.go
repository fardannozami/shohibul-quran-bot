package motivation

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Engine struct {
	ayat   []string
	hadith []string
	client *http.Client
}

func NewEngine() *Engine {
	return &Engine{
		ayat: []string{
			"Allah berfirman: “Ingatlah, hanya dengan mengingat Allah-lah hati menjadi tenteram.” (QS. Ar-Ra'd: 28)",
			"Allah berfirman: “Dan sesungguhnya telah Kami mudahkan Al-Qur'an untuk pelajaran, maka adakah orang yang mengambil pelajaran?” (QS. Al-Qamar: 17)",
			"Allah berfirman: “Bacalah dengan (menyebut) nama Tuhanmu yang menciptakan.” (QS. Al-Alaq: 1)",
		},
		hadith: []string{
			"“Sebaik-baik kalian adalah yang mempelajari Al-Qur'an dan mengajarkannya.” (HR. Bukhari)",
			"“Bacalah Al-Qur'an, karena sesungguhnya ia akan datang pada hari kiamat sebagai pemberi syafaat bagi pembacanya.” (HR. Muslim)",
			"“Orang yang lancar membaca Al-Qur'an akan bersama para malaikat yang mulia dan senantiasa taat...” (HR. Bukhari & Muslim)",
			"“Barangsiapa membaca satu huruf dari Kitabullah, maka baginya satu kebaikan...” (HR. Tirmidzi)",
			"“Sesungguhnya Allah mengangkat derajat sebuah kaum dengan kitab ini (Al-Qur'an), dan merendahkan kaum yang lain dengannya.” (HR. Muslim)",
		},
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetRandomMotivation picks a random quote from either list
func (e *Engine) GetRandomMotivation() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(2) == 0 {
		return e.GetRandomAyat()
	}
	return e.GetRandomHadith()
}

func (e *Engine) GetRandomAyat() string {
	// Try API first
	ayat, err := e.fetchRandomAyat()
	if err == nil && ayat != "" {
		return ayat
	}

	// Fallback to hard-coded list
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return e.ayat[r.Intn(len(e.ayat))]
}

func (e *Engine) GetRandomHadith() string {
	// Try API first
	hadist, err := e.fetchRandomHadith()
	if err == nil && hadist != "" {
		return hadist
	}

	// Fallback to hard-coded list
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return e.hadith[r.Intn(len(e.hadith))]
}

func (e *Engine) fetchRandomAyat() (string, error) {
	// API v4: /verses/random with translation ID 33 (Indonesian Kemenag)
	url := "https://api.quran.com/api/v4/verses/random?translations=33"
	resp, err := e.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var result struct {
		Verse struct {
			VerseKey     string `json:"verse_key"`
			Translations []struct {
				Text string `json:"text"`
			} `json:"translations"`
		} `json:"verse"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if len(result.Verse.Translations) > 0 {
		cleanText := result.Verse.Translations[0].Text
		// Remove HTML tags if any (basic cleaning)
		cleanText = strings.ReplaceAll(cleanText, "<sup", " <sup") // Add space before sup
		// We could use a more robust tag remover, but for simple quotes this might be enough
		return fmt.Sprintf("Allah berfirman: \"%s\" (QS. %s)", cleanText, result.Verse.VerseKey), nil
	}

	return "", fmt.Errorf("no translation found")
}

func (e *Engine) fetchRandomHadith() (string, error) {
	// Using api.hadith.gading.dev for random Bukhari
	// Range is typically 1 - 7008 for Bukhari in this API
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hadithNum := r.Intn(7008) + 1
	url := fmt.Sprintf("https://api.hadith.gading.dev/books/bukhari/%d", hadithNum)

	resp, err := e.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Id       string `json:"id"`
			Contents struct {
				Id string `json:"id"`
			} `json:"contents"`
		} `json:"data"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.Data.Contents.Id != "" {
		return fmt.Sprintf("\"%s\" (HR. Bukhari no. %d)", result.Data.Contents.Id, hadithNum), nil
	}

	return "", fmt.Errorf("no hadith content found")
}
