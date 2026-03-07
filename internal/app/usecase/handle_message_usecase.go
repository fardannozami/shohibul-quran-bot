package usecase

import (
	"context"
	"fmt"
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
	isReport, pages := uc.parser.Parse(msg)
	if isReport {
		return uc.gameEngine.ProcessReport(ctx, userID, name, pages, message)
	}

	// 2. Handle simple commands
	if strings.Contains(msg, "#leaderboard") {
		return uc.handleLeaderboard(ctx)
	}

	if strings.Contains(msg, "#mystats") {
		return uc.handleMyStats(ctx, userID, name)
	}

	if strings.Contains(msg, "#achievements") {
		return uc.handleAchievements(ctx)
	}

	return "", nil
}

func (uc *HandleMessageUsecase) handleLeaderboard(ctx context.Context) (string, error) {
	users, err := uc.repo.GetAllUsers(ctx)
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "Belum ada data member.", nil
	}

	resp := "🏆 *LEADERBOARD SHOHIBUL QURAN* 🏆\n\n"
	
	for i, u := range users {
		if i >= 10 { // top 10 only
			break
		}
		resp += fmt.Sprintf("%d. %s - Lvl %d | %d XP | Streak: %d🔥\n", i+1, u.Name, u.Level, u.XP, u.Streak)
	}

	return resp, nil
}

func (uc *HandleMessageUsecase) handleMyStats(ctx context.Context, userID, name string) (string, error) {
	user, err := uc.repo.GetUser(ctx, userID)
	if err != nil || user == nil {
		return fmt.Sprintf("Maaf %s, kamu belum ada data. Ayo lapor bacaan dulu!", name), nil
	}

	badges, _ := uc.repo.GetBadgesByUser(ctx, userID)
	
	resp := fmt.Sprintf("📊 *STATISTIK %s* 📊\n\n", name)
	resp += fmt.Sprintf("Level: %d\n", user.Level)
	resp += fmt.Sprintf("Total XP: %d\n", user.XP)
	resp += fmt.Sprintf("Current Streak: %d🔥\n", user.Streak)
	
	if len(badges) > 0 {
		resp += "\n🏅 *Badges:*\n"
		for _, b := range badges {
			resp += fmt.Sprintf("- %s\n", b.Badge)
		}
	} else {
		resp += "\nBelum ada badge. Terus semangat tilawah!"
	}

	return resp, nil
}

func (uc *HandleMessageUsecase) handleAchievements(ctx context.Context) (string, error) {
	resp := "🏅 *AVAILABLE ACHIEVEMENTS* 🏅\n\n"
	resp += "- 🏅 First Blood: First time reporting\n"
	resp += "- 🔥 Streak Master 7: 7 Days Streak\n"
	resp += "- 🔥 Streak Master 30: 30 Days Streak\n"
	resp += "- 📖 Bookworm 1 Juz: Read 1 Juz in a single day\n"

	return resp, nil
}
