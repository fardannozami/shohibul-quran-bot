package tafsir

import (
	"context"
	"strings"
	"testing"
)

func TestSummarizeFetched(t *testing.T) {
	e := NewEngine()
	ctx := context.Background()

	got, err := e.Summarize(ctx, 1, 1)
	if err != nil {
		t.Fatalf("Summarize(1,1) error: %v", err)
	}
	if got == "" {
		t.Fatal("Summarize(1,1) returned empty")
	}
	t.Logf("Al-Fatihah:1 => %s", got)
	if strings.Contains(got, "<") {
		t.Errorf("summary still contains HTML tags: %s", got)
	}

	got2, err := e.Summarize(ctx, 2, 255)
	if err != nil {
		t.Fatalf("Summarize(2,255) error: %v", err)
	}
	if got2 == "" {
		t.Fatal("Summarize(2,255) returned empty")
	}
	t.Logf("Al-Baqarah:255 => %s", got2)
}

func TestSummarizeTruncates(t *testing.T) {
	// A single very long sentence (no internal breaks) must be capped with an ellipsis.
	long := strings.Repeat("Ini kata yang sangat panjang ", 30) + "berakhir mulai sekarang."
	got := summarize(long)
	n := len([]rune(got))
	if n > maxSummaryRunes+3 {
		t.Errorf("summary too long: %d runes", n)
	}
	if !strings.HasSuffix(got, "…") {
		t.Errorf("expected trailing ellipsis, got %q", got)
	}
}

func TestSummarizeRangeFetched(t *testing.T) {
	e := NewEngine()
	ctx := context.Background()

	// Whole-surah report (As-Sajdah 1-30) should produce a summary spanning
	// more than just verse 1.
	single, err := e.Summarize(ctx, 32, 1)
	if err != nil {
		t.Fatalf("Summarize(32,1) error: %v", err)
	}
	whole, err := e.SummarizeRange(ctx, 32, 1, 30)
	if err != nil {
		t.Fatalf("SummarizeRange(32,1,30) error: %v", err)
	}
	if whole == "" {
		t.Fatal("SummarizeRange(32,1,30) returned empty")
	}
	t.Logf("single: %s", single)
	t.Logf("whole: %s", whole)
	if len([]rune(whole)) < len([]rune(single)) {
		t.Errorf("range summary %d runes should be at least as long as single %d", len([]rune(whole)), len([]rune(single)))
	}

	// A single-ayah range should behave like Summarize.
	one, err := e.SummarizeRange(ctx, 32, 1, 1)
	if err != nil {
		t.Fatalf("SummarizeRange(32,1,1) error: %v", err)
	}
	if one != single {
		t.Errorf("single-ayah range mismatch:\n got  %q\n want %q", one, single)
	}
}

type fakeAI struct {
	summary string
	err     error
	gotName string
	gotLabel string
	gotText  string
}

func (f *fakeAI) GenerateTafsirSummary(_ context.Context, surahName, ayahLabel, teks string) (string, error) {
	f.gotName = surahName
	f.gotLabel = ayahLabel
	f.gotText = teks
	return f.summary, f.err
}

func TestSummarizeUsesAI(t *testing.T) {
	ctx := context.Background()

	t.Run("single ayah", func(t *testing.T) {
		fake := &fakeAI{summary: "Ringkasan AI ayat 1."}
		e := NewEngine()
		e.SetAI(fake)

		got, err := e.Summarize(ctx, 32, 1)
		if err != nil {
			t.Fatalf("Summarize error: %v", err)
		}
		if got != "Ringkasan AI ayat 1." {
			t.Errorf("expected AI summary, got %q", got)
		}
		if fake.gotLabel != "1" {
			t.Errorf("expected label 1, got %q", fake.gotLabel)
		}
		if !strings.Contains(fake.gotText, "Alif Lam Mim") {
			t.Errorf("expected tafsir text sent to AI, got %q", fake.gotText)
		}
	})

	t.Run("range uses AI when available", func(t *testing.T) {
		fake := &fakeAI{summary: "Ringkasan AI 1-30."}
		e := NewEngine()
		e.SetAI(fake)

		got, err := e.SummarizeRange(ctx, 32, 1, 30)
		if err != nil {
			t.Fatalf("SummarizeRange error: %v", err)
		}
		if got != "Ringkasan AI 1-30." {
			t.Errorf("expected AI range summary, got %q", got)
		}
		if fake.gotLabel != "1-30" {
			t.Errorf("expected label 1-30, got %q", fake.gotLabel)
		}
		if fake.gotName != "As-Sajdah" {
			t.Errorf("expected surah name As-Sajdah, got %q", fake.gotName)
		}
	})

	t.Run("falls back when AI errors", func(t *testing.T) {
		fake := &fakeAI{err: context.DeadlineExceeded}
		e := NewEngine()
		e.SetAI(fake)

		got, err := e.Summarize(ctx, 32, 1)
		if err != nil {
			t.Fatalf("Summarize error: %v", err)
		}
		if got == "" {
			t.Fatal("expected heuristic fallback, got empty")
		}
	})
}