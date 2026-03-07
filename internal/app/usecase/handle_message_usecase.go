package usecase

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/fardannozami/shohibul-quran-bot/internal/app/gamification"
	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"github.com/fardannozami/shohibul-quran-bot/internal/parser"
)

type HandleMessageUsecase struct {
	repo       domain.BotRepository
	parser     *parser.ReportParser
	gameEngine *gamification.Engine
}

func NewHandleMessageUsecase(repo domain.BotRepository, parser *parser.ReportParser, gameEngine *gamification.Engine) *HandleMessageUsecase {
	return &HandleMessageUsecase{
		repo:       repo,
		parser:     parser,
		gameEngine: gameEngine,
	}
}

func (uc *HandleMessageUsecase) Execute(ctx context.Context, userID, name, message string) (string, error) {
	msg := strings.ToLower(strings.TrimSpace(message))

	// 1. Check if it's a report (contains alhamdulillah)
	result := uc.parser.Parse(msg)
	if result.IsReport {
		return uc.gameEngine.ProcessReport(ctx, userID, name, result, message)
	}

	// 2. Handle simple commands (support both # and ! prefix)
	if strings.Contains(msg, "#leaderboard") || strings.Contains(msg, "!leaderboard") {
		return uc.handleLeaderboard(ctx)
	}
	if strings.Contains(msg, "#mystats") || strings.Contains(msg, "!stats") {
		return uc.handleMyStats(ctx, userID, name)
	}
	if strings.Contains(msg, "#achievements") || strings.Contains(msg, "!achievements") {
		return uc.handleAchievements(ctx)
	}
	if strings.Contains(msg, "!target") {
		return uc.handleTarget(ctx)
	}

	return "", nil
}

func (uc *HandleMessageUsecase) handleTarget(ctx context.Context) (string, error) {
	// For now, return a fixed goal message based on PRD suggestions
	// In the future, this can be calculated from weekly reports
	resp := "🎯 *Target Komunitas Shohibul Qur'an*\n\n"
	resp += "Bersama kita targetkan khatam *50 Juz* setiap minggunya! 🔥\n\n"
	resp += "Progress minggu ini: _(Segera hadir fitur akumulasi mingguan!)_\n\n"
	resp += "Semangat tilawah semuanya! 📖✨"
	return resp, nil
}

func (uc *HandleMessageUsecase) handleLeaderboard(ctx context.Context) (string, error) {
	users, err := uc.repo.GetAllUsers(ctx)
	if err != nil {
		return "", err
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i].XP > users[j].XP
	})

	if len(users) == 0 {
		return "Belum ada data anggota. Ayo mulai tilawah dan laporkan bacaan pertamamu! 📖", nil
	}

	resp := "📖 *PERINGKAT SHOHIBUL QUR'AN* 📖\n"
	resp += "━━━━━━━━━━━━━━━\n\n"

	medals := []string{"🥇", "🥈", "🥉"}
	for i, u := range users {
		if i >= 10 {
			break
		}
		medal := fmt.Sprintf("%d.", i+1)
		if i < 3 {
			medal = medals[i]
		}
		resp += fmt.Sprintf("%s *%s*\n    Level %d  •  %d XP  •  %d hari 🔥\n\n", medal, u.Name, u.Level, u.XP, u.Streak)
	}

	resp += "━━━━━━━━━━━━━━━\n"
	resp += "بارك الله فيكم\nSemoga Allah memberkahi kita semua 🤲"

	return resp, nil
}

func (uc *HandleMessageUsecase) handleMyStats(ctx context.Context, userID, name string) (string, error) {
	user, err := uc.repo.GetUser(ctx, userID)
	if err != nil || user == nil {
		return fmt.Sprintf("Afwan %s, belum ada data tilawahmu. Yuk mulai baca Al-Qur'an dan laporkan! 📖🤲", name), nil
	}

	badges, _ := uc.repo.GetBadgesByUser(ctx, userID)
	
	resp := fmt.Sprintf("📊 *Statistik Tilawah*\n━━━━━━━━━━━━━━━\n\n👤  *%s*\n\n", name)
	resp += fmt.Sprintf("🕌  Level: *%d*\n", user.Level)
	resp += fmt.Sprintf("⭐  Total XP: *%d*\n", user.XP)
	resp += fmt.Sprintf("🔥  Istiqomah: *%d hari*\n", user.Streak)

	if len(badges) > 0 {
		resp += "\n━━━━━━━━━━━━━━━\n"
		resp += "✨ *Pencapaian:*\n\n"
		for _, b := range badges {
			resp += fmt.Sprintf("  %s\n", b.Badge)
		}
	} else {
		resp += "\n━━━━━━━━━━━━━━━\n"
		resp += "Belum ada pencapaian.\nTerus istiqomah, insyaAllah segera didapat! 🤲"
	}

	resp += "\n━━━━━━━━━━━━━━━\n"
	resp += "Semoga Allah memudahkan tilawahmu 🤍"

	return resp, nil
}

func (uc *HandleMessageUsecase) handleAchievements(ctx context.Context) (string, error) {
	resp := "✨ *PENCAPAIAN SHOHIBUL QUR'AN* ✨\n"
	resp += "━━━━━━━━━━━━━━━\n\n"
	resp += "🕌  *Langkah Pertama*\n      Tilawah pertamamu tercatat\n\n"
	resp += "🔥  *Sahabat Qur'an*\n      Istiqomah 7 hari berturut-turut\n\n"
	resp += "🌙  *Ahlul Qur'an*\n      Istiqomah 30 hari berturut-turut\n\n"
	resp += "📖  *Khatam Juz*\n      Membaca 1 juz (20 halaman) dalam sehari\n\n"
	resp += "━━━━━━━━━━━━━━━\n"
	resp += "Terus istiqomah!\nSetiap huruf yang dibaca bernilai kebaikan 🤲"

	return resp, nil
}
