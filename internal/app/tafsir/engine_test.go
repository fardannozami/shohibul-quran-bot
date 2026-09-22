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