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

const tafsirSystemPrompt = `Kamu adalah penulis ringkasan tafsir Al-Qur'an dalam bahasa Indonesia yang akurat dan mudah dipahami.

Tugas & aturan:
- Buat ringkasan dari teks tafsir Kemenag yang diberikan, berdasarkan isi teks tersebut saja, jangan menambah informasi dari luar.
- Tulis 3-5 kalimat paragraf yang mengalir; jangan pakai daftar atau bullet.
- Ungkapkan intisarinya dengan bahasa yang enak dibaca, bukan menempel teks asli mentah-mentah.
- Hindari pengulangan pembuka seperti "Ayat ini menerangkan bahwa..." di awal setiap kalimat.
- Jangan gunakan emoji atau format markdown.`

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
	return e.generate(ctx, systemPrompt, userMessage, 0.7, 800, -1)
}

// GenerateTafsirSummary asks the model to produce a concise Indonesian summary
// of the given tafsir text (Kemenag via equran.id). Thinking is disabled so
// the output never gets cut off by the internal thinking budget.
func (e *Engine) GenerateTafsirSummary(ctx context.Context, surahName, ayahLabel, teks string) (string, error) {
	prompt := fmt.Sprintf("Ringkas tafsir QS. %s:%s berikut dalam 3-5 kalimat yang enak dibaca:\n\n%s", surahName, ayahLabel, teks)
	return e.generate(ctx, tafsirSystemPrompt, prompt, 0.3, 1000, 0)
}

func (e *Engine) generate(ctx context.Context, system, user string, temperature float64, maxTokens, thinkingTokens int) (string, error) {
	if e.apiKey == "" {
		return "", fmt.Errorf("gemini api key is empty")
	}
	url := geminiEndpoint + "?key=" + e.apiKey

	genConfig := map[string]interface{}{
		"temperature":     temperature,
		"maxOutputTokens": maxTokens,
	}
	if thinkingTokens >= 0 {
		genConfig["thinkingConfig"] = map[string]interface{}{
			"thinkingBudget": thinkingTokens,
		}
	}

	payload := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{{"text": system}},
		},
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": user}}},
		},
		"generationConfig": genConfig,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	var respBody []byte
	var status int
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := e.client.Do(req)
		if err != nil {
			if attempt == 2 {
				return "", err
			}
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Duration(attempt+1) * time.Second):
			}
			continue
		}
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		status = resp.StatusCode
		if status == http.StatusOK {
			break
		}
		if status == 429 || status == 503 || status >= 500 {
			if attempt == 2 {
				return "", fmt.Errorf("gemini API error: %d %s", status, strings.TrimSpace(string(respBody)))
			}
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Duration(attempt+1)*2 * time.Second):
			}
			continue
		}
		return "", fmt.Errorf("gemini API error: %d %s", status, strings.TrimSpace(string(respBody)))
	}
	if status != http.StatusOK {
		return "", fmt.Errorf("gemini API error: %d %s", status, strings.TrimSpace(string(respBody)))
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
