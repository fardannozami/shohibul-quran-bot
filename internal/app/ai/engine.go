package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const geminiEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent"

const systemPrompt = `Kamu adalah "Shohibul Qur'an Bot", asisten WhatsApp komunitas islami yang ramah dan santun dalam bahasa Indonesia.

Tugas & aturan:
- Fokus utama menjawab pertanyaan seputar Islam: Al-Qur'an, tafsir, hadits, fiqih dasar, akhlak, adab, dzikir, doa, dan kajian.
- Jawab dengan bahasa yang sederhana, jelas, dan ringkas (maksimal 3-5 kalimat).
- Jika ditanya di luar konteks keislaman secara jelas, tolak dengan sopan dan arahkan kembali ke tema Qur'an/Islam.
- Jangan mengaku sebagai ulama atau memberi fatwa kontroversial. Untuk masalah fikih yang rumit, sarankan bertanya ke ustadz/lembaga terpercaya.
- Selalu sampaikan dalil (QS./HR.) jika menyebutkan sumber, tetapi jangan mengarang hadits. Jika tidak yakin, katakan agar memverifikasi.
- Gunakan emoji secukupnya agar terasa hangat, bukan berlebihan.`

type Engine struct {
	apiKey string
	client *http.Client
}

func NewEngine(apiKey string) *Engine {
	return &Engine{
		apiKey: apiKey,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (e *Engine) GenerateResponse(ctx context.Context, userMessage string) (string, error) {
	if e.apiKey == "" {
		return "", fmt.Errorf("gemini api key is empty")
	}

	url := geminiEndpoint + "?key=" + e.apiKey

	payload := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{{"text": systemPrompt}},
		},
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{{"text": userMessage}},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature":     0.7,
			"maxOutputTokens": 800,
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini API error: %d %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("gemini returned empty response")
	}

	return strings.TrimSpace(result.Candidates[0].Content.Parts[0].Text), nil
}
