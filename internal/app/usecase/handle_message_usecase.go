package usecase

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/app/gamification"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/motivation"
	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"github.com/fardannozami/shohibul-quran-bot/internal/parser"
)

type HandleMessageUsecase struct {
	repo       domain.BotRepository
	parser     *parser.ReportParser
	gameEngine *gamification.Engine
	motEngine  *motivation.Engine
}

func NewHandleMessageUsecase(repo domain.BotRepository, parser *parser.ReportParser, gameEngine *gamification.Engine, motEngine *motivation.Engine) *HandleMessageUsecase {
	return &HandleMessageUsecase{
		repo:       repo,
		parser:     parser,
		gameEngine: gameEngine,
		motEngine:  motEngine,
	}
}

func (uc *HandleMessageUsecase) Execute(ctx context.Context, userID, name, message string, groupID string) (string, error) {
	msg := strings.ToLower(strings.TrimSpace(message))

	// 1. Check if it's a report (contains alhamdulillah)
	results := uc.parser.Parse(msg)
	if len(results) > 0 {
		return uc.gameEngine.ProcessReports(ctx, userID, name, results, message, groupID)
	}

	// 2. Handle simple commands (support both # and ! prefix)
	if strings.Contains(msg, "#leaderboard") || strings.Contains(msg, "!leaderboard") {
		return uc.handleLeaderboard(ctx, groupID)
	}
	if strings.Contains(msg, "#mystats") || strings.Contains(msg, "!stats") {
		return uc.handleMyStats(ctx, userID, name, groupID)
	}
	if strings.Contains(msg, "#achievements") || strings.Contains(msg, "!achievements") {
		return uc.handleAchievements(ctx)
	}
	if (strings.HasPrefix(msg, "#settarget") || strings.HasPrefix(msg, "!settarget")) && len(msg) > 10 {
		return uc.handleSetTarget(ctx, userID, name, groupID, message)
	}
	if strings.Contains(msg, "!target") {
		return uc.handleTarget(ctx, groupID)
	}
	if strings.Contains(msg, "!ayat") {
		return "📖 *Ayat Qur'an*\n\n" + uc.motEngine.GetRandomAyat(), nil
	}
	if strings.Contains(msg, "!hadith") || strings.Contains(msg, "!hadist") {
		return "🌙 *Hadist*\n\n" + uc.motEngine.GetRandomHadith(), nil
	}

	if strings.Contains(msg, "!cara") || strings.Contains(msg, "#cara") {
		return uc.handleHelp(), nil
	}

	return "", nil
}

func (uc *HandleMessageUsecase) handleHelp() string {
	resp := "📖 *CARA PENGGUNAAN BOT SHOHIBUL QUR'AN* 📖\n\n"
	resp += "📌 *Cara Laporan Bacaan:*\n"
	resp += "Cukup kirim pesan mengandung kata \"Alhamdulillah\" dan jumlah bacaanmu.\n"
	resp += "Contoh:\n"
	resp += "- Alhamdulillah sudah baca 2 halaman\n"
	resp += "- Alhamdulillah 1 juz\n"
	resp += "- Alhamdulillah Al-Baqarah 1-30\n\n"
	resp += "🎯 *Cara Atur Target Harian:*\n"
	resp += "Gunakan perintah `!settarget` diikuti jumlah dan satuan (`juz` atau `halaman`).\n"
	resp += "Contoh:\n"
	resp += "- `!settarget 1 juz`\n"
	resp += "- `!settarget 10 halaman`\n"
	resp += "- `!settarget 0` (untuk menghapus target)\n\n"
	resp += "📊 *Perintah Lain:*\n"
	resp += "- `!stats` atau `#mystats`: Lihat statistik pribadimu\n"
	resp += "- `!leaderboard`: Lihat peringkat 10 besar\n"
	resp += "- `!target`: Lihat progress target komunitas\n"
	resp += "- `!achievements`: Daftar pencapaian yang bisa diraih\n\n"
	resp += "Semoga bermanfaat! Barakallahu fiikum. 🤍"
	return resp
}

func (uc *HandleMessageUsecase) handleTarget(ctx context.Context, groupID string) (string, error) {
	users, err := uc.repo.GetAllUsers(ctx, groupID)
	if err != nil {
		return "", err
	}
	memberCount := len(users)
	if memberCount == 0 {
		memberCount = 1 // fallback
	}

	// Target: members * 7 juz
	targetJuz := memberCount * 7
	targetPages := targetJuz * 20 // 1 juz = 20 pages

	// Calculate weekly progress (Monday to now)
	now := time.Now()
	// Monday is 1, Sunday is 0 in Time.Weekday()? No, Sunday=0, Monday=1, ... Saturday=6
	daysSinceMonday := int(now.Weekday()) - 1
	if daysSinceMonday < 0 {
		daysSinceMonday = 6 // Sunday
	}
	startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -daysSinceMonday)

	currentPages, err := uc.repo.GetTotalPagesInRange(ctx, startOfWeek, now, groupID)
	if err != nil {
		currentPages = 0
	}

	currentJuz := float64(currentPages) / 20.0

	resp := "🎯 *Target Komunitas Shohibul Qur'an*\n\n"
	resp += fmt.Sprintf("Anggota aktif: *%d orang*\n", memberCount)
	resp += fmt.Sprintf("Target mingguan: *%d Juz* (%d halaman) 🔥\n\n", targetJuz, targetPages)
	
	progressBar := uc.generateProgressBar(currentPages, targetPages)
	resp += fmt.Sprintf("📊 *Progress Minggu Ini:*\n%s\n", progressBar)
	resp += fmt.Sprintf("Tercapai: *%.1f Juz* (%d halaman)\n\n", currentJuz, currentPages)
	
	if currentPages >= targetPages {
		resp += "🎉 *MABRUK!* Target mingguan tercapai. Teruslah istiqomah! 🚀"
	} else {
		resp += "Semangat tilawah semuanya! Sedikit lagi target tercapai 📖✨"
	}
	
	return resp, nil
}

func (uc *HandleMessageUsecase) generateProgressBar(current, target int) string {
	if target <= 0 {
		return "[░░░░░░░░░░] 0%"
	}
	percent := (current * 100) / target
	if percent > 100 {
		percent = 100
	}
	
	filled := percent / 10
	bar := ""
	for i := 0; i < 10; i++ {
		if i < filled {
			bar += "▓"
		} else {
			bar += "░"
		}
	}
	return fmt.Sprintf("[%s] %d%%", bar, percent)
}


func (uc *HandleMessageUsecase) handleLeaderboard(ctx context.Context, groupID string) (string, error) {
	users, err := uc.repo.GetAllUsers(ctx, groupID)
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

func (uc *HandleMessageUsecase) handleMyStats(ctx context.Context, userID, name string, groupID string) (string, error) {
	user, err := uc.repo.GetUser(ctx, userID, groupID)
	if err != nil || user == nil {
		return fmt.Sprintf("Afwan %s, belum ada data tilawahmu di grup ini. Yuk mulai baca Al-Qur'an dan laporkan! 📖🤲", name), nil
	}

	badges, _ := uc.repo.GetBadgesByUser(ctx, userID, groupID)
	
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

func (uc *HandleMessageUsecase) handleSetTarget(ctx context.Context, userID, name, groupID, message string) (string, error) {
	// message format: !settarget 1 juz OR !settarget 5 halaman OR !settarget 0
	msg := strings.TrimSpace(message)
	parts := strings.Fields(msg)
	if len(parts) < 2 {
		return "Format salah. Contoh: `!settarget 1 juz` atau `!settarget 5 halaman` atau `!settarget 0` (untuk menghapus)", nil
	}

	valStr := parts[1]
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return "Nilai target harus berupa angka. Contoh: `!settarget 20`", nil
	}

	if val < 0 {
		return "Nilai target tidak boleh negatif.", nil
	}

	pages := val
	unit := ""
	if len(parts) >= 3 {
		unit = strings.ToLower(parts[2])
	}

	if unit == "juz" {
		pages = val * 20
	} else if unit == "halaman" || unit == "hal" || unit == "hlm" || unit == "" {
		pages = val
	} else {
		return "Satuan tidak dikenal. Gunakan `juz` atau `halaman`.", nil
	}

	user, err := uc.repo.GetUser(ctx, userID, groupID)
	if err != nil {
		return "", err
	}

	if user == nil {
		user = &domain.User{
			ID:          userID,
			GroupID:     groupID,
			Phone:       uc.repo.ResolveLIDToPhone(ctx, userID),
			Name:        name,
			XP:          0,
			Level:       1,
			Streak:      0,
			JoinedAt:    time.Now(),
			DailyTarget: pages,
		}
		if err := uc.repo.CreateUser(ctx, user); err != nil {
			return "", err
		}
	} else {
		user.DailyTarget = pages
		if err := uc.repo.UpdateUser(ctx, user); err != nil {
			return "", err
		}
	}

	if pages == 0 {
		return fmt.Sprintf("✅ *%s*, target harianmu telah dihapus.", name), nil
	}

	targetDesc := fmt.Sprintf("%d halaman", pages)
	if unit == "juz" {
		targetDesc = fmt.Sprintf("%d Juz (%d halaman)", val, pages)
	}

	return fmt.Sprintf("✅ *%s*, target harianmu berhasil diatur: *%s*.\n\nSemoga Allah memudahkan tilawahmu! 🤲", name, targetDesc), nil
}
