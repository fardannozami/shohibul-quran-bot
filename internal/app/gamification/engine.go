package gamification

import (
	"context"
	"fmt"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
)

type Engine struct {
	repo domain.BotRepository
}

func NewEngine(repo domain.BotRepository) *Engine {
	return &Engine{repo: repo}
}

// ProcessReport computes XP, streaks, and badges for an incoming report.
// Returns a structured message string for reporting back to the user.
func (e *Engine) ProcessReport(ctx context.Context, userID, name string, pages int, messageText string) (string, error) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	yesterday := today.AddDate(0, 0, -1)

	// Ensure User exists
	user, err := e.repo.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}

	if user == nil {
		user = &domain.User{
			ID:       userID,
			Phone:    e.repo.ResolveLIDToPhone(ctx, userID), // simple mapping
			Name:     name,
			XP:       0,
			Level:    1,
			Streak:   0,
			JoinedAt: now,
		}
		if err := e.repo.CreateUser(ctx, user); err != nil {
			return "", err
		}
	} else if user.Name != name {
		user.Name = name
		_ = e.repo.UpdateUser(ctx, user)
	}

	// Fetch today's progress to check if they already reported today
	todayProgress, err := e.repo.GetDailyProgress(ctx, userID, today)
	if err != nil {
		return "", err
	}

	isNewStreak := false
	if todayProgress == nil || todayProgress.ReportsCount == 0 {
		// New report for today. Let's check yesterday's progress for streak logic.
		yesterdayProgress, err := e.repo.GetDailyProgress(ctx, userID, yesterday)
		if err != nil {
			return "", err
		}

		if yesterdayProgress != nil && yesterdayProgress.Pages > 0 {
			// Streak continues
			user.Streak += 1
		} else {
			// Streak resets
			user.Streak = 1
			isNewStreak = true
		}
	}

	// Calculate XP: let's say 1 page = 10 XP
	xpGained := pages * 10
	user.XP += xpGained

	// Calculate Level: very simple 1 level = every 100 XP
	oldLevel := user.Level
	user.Level = (user.XP / 100) + 1

	// Save User
	if err := e.repo.UpdateUser(ctx, user); err != nil {
		return "", err
	}

	// Update Daily Progress
	if todayProgress == nil {
		todayProgress = &domain.DailyProgress{
			UserID:       userID,
			Date:         today,
			Pages:        pages,
			ReportsCount: 1,
		}
	} else {
		todayProgress.Pages += pages
		todayProgress.ReportsCount += 1
	}

	if err := e.repo.UpsertDailyProgress(ctx, todayProgress); err != nil {
		return "", err
	}

	// Insert Report Log
	reportLog := &domain.ReportLog{
		ID:      fmt.Sprintf("%s-%d", userID, now.UnixNano()),
		UserID:  userID,
		Pages:   pages,
		Message: messageText,
		Date:    now,
	}
	if err := e.repo.InsertReport(ctx, reportLog); err != nil {
		return "", err
	}

	// Check and Grant Badges
	badgeMsg := e.checkBadges(ctx, user, todayProgress)

	// Format response message
	resp := fmt.Sprintf("Laporan diterima, %s sudah tilawah %d halaman hari ini. Lanjutkan 🔥\n", name, todayProgress.Pages)
	
	if isNewStreak && user.Streak == 1 {
		resp += "Semangat memulai streak baru!\n"
	} else {
		resp += fmt.Sprintf("(🔥 Streak: %d hari)\n", user.Streak)
	}

	resp += fmt.Sprintf("⭐ +%d XP (Total XP: %d)", xpGained, user.XP)

	if user.Level > oldLevel {
		resp += fmt.Sprintf("\n\n🎉 LEVEL UP! Sekarang kamu Level %d! 🎉", user.Level)
	}

	if badgeMsg != "" {
		resp += "\n" + badgeMsg
	}

	return resp, nil
}

func (e *Engine) checkBadges(ctx context.Context, user *domain.User, todayProgress *domain.DailyProgress) string {
	var newBadges []string

	checkAndGrant := func(badgeName string, condition bool) {
		if !condition {
			return
		}
		// check if already has badge
		existing, _ := e.repo.GetBadgesByUser(ctx, user.ID)
		for _, b := range existing {
			if b.Badge == badgeName {
				return
			}
		}

		// grant badge
		newBadge := &domain.BadgeLog{
			UserID:    user.ID,
			Badge:     badgeName,
			CreatedAt: time.Now(),
		}
		_ = e.repo.InsertBadge(ctx, newBadge)
		newBadges = append(newBadges, badgeName)
	}

	// Definitions
	checkAndGrant("🏅 First Blood", user.XP > 0)
	checkAndGrant("🔥 Streak Master 7", user.Streak >= 7)
	checkAndGrant("🔥 Streak Master 30", user.Streak >= 30)
	checkAndGrant("📖 Bookworm 1 Juz", todayProgress.Pages >= 20)

	msg := ""
	if len(newBadges) > 0 {
		msg += "\n🎉 ACHIEVEMENT UNLOCKED! 🎉"
		for _, b := range newBadges {
			msg += "\n" + b
		}
	}
	return msg
}
