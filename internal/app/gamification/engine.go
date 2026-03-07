package gamification

import (
	"context"
	"fmt"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"github.com/fardannozami/shohibul-quran-bot/internal/parser"
)

type Engine struct {
	repo domain.BotRepository
}

func NewEngine(repo domain.BotRepository) *Engine {
	return &Engine{repo: repo}
}

// ProcessReport computes XP, streaks, and badges for an incoming report.
// Returns a structured message string for reporting back to the user.
func (e *Engine) ProcessReport(ctx context.Context, userID, name string, result parser.ParseResult, messageText string) (string, error) {
	pages := result.Pages
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

	// Calculate XP: 10 base for reporting + 2 per page
	xpGained := 10 + (pages * 2)
	streakBonus := 0

	// If this is the first report of the day and they hit a 7-day milestone, give bonus
	if (todayProgress == nil || todayProgress.ReportsCount == 0) && user.Streak > 0 && user.Streak%7 == 0 {
		streakBonus = 20
		xpGained += streakBonus
	}

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
	resp := "بسم الله الرحمن الرحيم\n"
	resp += "━━━━━━━━━━━━━━━\n\n"

	// Show surah detail if available
	if surahInfo := result.FormatSurahInfo(); surahInfo != "" {
		resp += fmt.Sprintf("MasyaAllah tabarakallah, %s telah tilawah\n\n", name)
		resp += fmt.Sprintf("%s (%d halaman)\n", surahInfo, pages)
		resp += fmt.Sprintf("Total tilawah hari ini: *%d halaman*\n\n", todayProgress.Pages)
	} else {
		resp += fmt.Sprintf("📖 MasyaAllah tabarakallah, %s telah tilawah *%d halaman* Al-Qur'an hari ini.\n\n", name, todayProgress.Pages)
	}

	if isNewStreak && user.Streak == 1 {
		resp += "🌱 Bismillah, semoga menjadi awal istiqomah yang barokah.\n\n"
	} else {
		resp += fmt.Sprintf("🔥 *Istiqomah %d hari* berturut-turut — MasyaAllah\n\n", user.Streak)
	}

	resp += "━━━━━━━━━━━━━━━\n"
	resp += fmt.Sprintf("⭐  +%d XP\n", xpGained)
	if streakBonus > 0 {
		resp += fmt.Sprintf("🎁  Bonus istiqomah +%d XP\n", streakBonus)
	}
	resp += fmt.Sprintf("📊  Total XP: *%d*\n", user.XP)
	resp += fmt.Sprintf("🕌  Level: *%d*\n", user.Level)
	resp += "━━━━━━━━━━━━━━━\n"

	if user.Level > oldLevel {
		resp += fmt.Sprintf("\n🌟 *Allahumma baarik*\nNaik ke *Level %d* Terus jaga tilawahmu 📖\n", user.Level)
	}

	if badgeMsg != "" {
		resp += badgeMsg
	}

	resp += "\nSemoga Allah memberkahi setiap huruf yang dibaca 🤲"

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
	checkAndGrant("🕌 Langkah Pertama — Bismillah, tilawah pertamamu tercatat!", user.XP > 0)
	checkAndGrant("🔥 Sahabat Qur'an — Istiqomah 7 hari berturut-turut", user.Streak >= 7)
	checkAndGrant("🌙 Ahlul Qur'an — Istiqomah 30 hari berturut-turut", user.Streak >= 30)
	checkAndGrant("📖 Pelajar Al-Qur'an — MasyaAllah, membaca 1 juz dalam sehari", todayProgress.Pages >= 20)
	checkAndGrant("📖 Hamalatul Qur'an — Luar biasa, membaca 2 juz dalam sehari", todayProgress.Pages >= 40)
	checkAndGrant("📖 Khadimul Qur'an — Tabarakallah, membaca 3 juz dalam sehari!", todayProgress.Pages >= 60)
	checkAndGrant("📖 Hafidzul Qur'an (Daily) — Allahu Akbar, membaca 5 juz dalam sehari!", todayProgress.Pages >= 100)

	msg := ""
	if len(newBadges) > 0 {
		msg += "\n\n✨ *Alhamdulillah, pencapaian baru!* ✨"
		for _, b := range newBadges {
			msg += "\n" + b
		}
		msg += "\n\nSemoga menjadi amal jariyah dan syafaat di hari akhir 🤲"
	}
	return msg
}
