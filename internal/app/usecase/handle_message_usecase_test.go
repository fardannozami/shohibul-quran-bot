package usecase

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
)

type mockRepo struct {
	domain.BotRepository
	user *domain.User
}

func (m *mockRepo) GetUser(ctx context.Context, id string, groupID string) (*domain.User, error) {
	return m.user, nil
}

func (m *mockRepo) CreateUser(ctx context.Context, user *domain.User) error {
	m.user = user
	return nil
}

func (m *mockRepo) UpdateUser(ctx context.Context, user *domain.User) error {
	m.user = user
	return nil
}

func (m *mockRepo) ResolveLIDToPhone(ctx context.Context, lid string) string {
	return lid
}

func (m *mockRepo) GetDailyProgress(ctx context.Context, userID string, groupID string, date time.Time) (*domain.DailyProgress, error) {
	return nil, nil
}

func (m *mockRepo) GetBadgesByUser(ctx context.Context, userID string, groupID string) ([]*domain.BadgeLog, error) {
	return nil, nil
}

func TestHandleSetTarget(t *testing.T) {
	repo := &mockRepo{}
	uc := &HandleMessageUsecase{
		repo: repo,
	}

	ctx := context.Background()
	userID := "user123"
	name := "Alice"
	groupID := "group456"

	t.Run("Set target pages", func(t *testing.T) {
		resp, err := uc.Execute(ctx, userID, name, "!settarget 10 halaman", groupID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !strings.Contains(resp, "10 halaman") {
			t.Errorf("expected response to contain 10 halaman, got %s", resp)
		}
		if repo.user.DailyTarget != 10 {
			t.Errorf("expected DailyTarget to be 10, got %d", repo.user.DailyTarget)
		}
	})

	t.Run("Set target juz", func(t *testing.T) {
		resp, err := uc.Execute(ctx, userID, name, "!settarget 1 juz", groupID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !strings.Contains(resp, "20 halaman") {
			t.Errorf("expected response to contain 20 halaman, got %s", resp)
		}
		if repo.user.DailyTarget != 20 {
			t.Errorf("expected DailyTarget to be 20, got %d", repo.user.DailyTarget)
		}
	})

	t.Run("Reset target", func(t *testing.T) {
		resp, err := uc.Execute(ctx, userID, name, "!settarget 0", groupID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !strings.Contains(resp, "dihapus") {
			t.Errorf("expected response to contain 'dihapus', got %s", resp)
		}
		if repo.user.DailyTarget != 0 {
			t.Errorf("expected DailyTarget to be 0, got %d", repo.user.DailyTarget)
		}
	})
	t.Run("Help command / !cara", func(t *testing.T) {
		resp, err := uc.Execute(ctx, userID, name, "!cara", groupID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !strings.Contains(resp, "CARA PENGGUNAAN") {
			t.Errorf("expected response to contain 'CARA PENGGUNAAN', got %s", resp)
		}
		if !strings.Contains(resp, "Single Surah") {
			t.Errorf("expected response to contain 'Single Surah', got %s", resp)
		}
		if !strings.Contains(resp, "Multiple Surah") {
			t.Errorf("expected response to contain 'Multiple Surah', got %s", resp)
		}
	})

	t.Run("Stats with target progress", func(t *testing.T) {
		// Set target first
		_, _ = uc.Execute(ctx, userID, name, "!settarget 10", groupID)
		
		resp, err := uc.Execute(ctx, userID, name, "!stats", groupID)
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !strings.Contains(resp, "🎯  Target: *10 hlm*") {
			t.Errorf("expected response to contain target info, got %s", resp)
		}
	})
}
