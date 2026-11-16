package scheduler

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/app/kajian"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/motivation"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/prayer"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/sunnah"
	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

// Nama slot notifikasi yang bisa dimatikan lewat ENV DISABLED_NOTIFS.
const (
	NotifPagi      = "pagi"      // 07:00 — ayat/hadist, kajian, countdown Ramadhan, pengingat dhuha
	NotifPengingat = "pengingat" // 18:00 — mention member yang belum laporan
	NotifMalam     = "malam"     // 21:00 — sunnah tidur / malam Jum'at
)

type CronService struct {
	client       *whatsmeow.Client
	repo         domain.BotRepository
	motEngine    *motivation.Engine
	prayerEngine *prayer.Engine
	kajianEngine *kajian.Engine
	sunnahEngine *sunnah.Engine
	groupIDs     []string
	disabled     map[string]bool
}

func NewCronService(client *whatsmeow.Client, repo domain.BotRepository, motEngine *motivation.Engine, prayerEngine *prayer.Engine, kajianEngine *kajian.Engine, sunnahEngine *sunnah.Engine, groupIDs []string, disabled map[string]bool) *CronService {
	return &CronService{
		client:       client,
		repo:         repo,
		motEngine:    motEngine,
		prayerEngine: prayerEngine,
		kajianEngine: kajianEngine,
		sunnahEngine: sunnahEngine,
		groupIDs:     groupIDs,
		disabled:     disabled,
	}
}

func (s *CronService) Start(ctx context.Context) {
	if len(s.groupIDs) == 0 {
		log.Println("No group IDs configured, skipping cron jobs")
		return
	}

	// Broadcast dibatasi 3 slot per hari agar grup tidak dibanjiri pesan bot:
	// pagi (07:00), pengingat laporan (18:00), dan malam (21:00).
	if s.enabled(NotifPagi) {
		go s.runDailyJob(ctx, 7, 0, "morning digest", s.executeMorningDigest)
	}
	if s.enabled(NotifPengingat) {
		go s.runDailyJob(ctx, 18, 0, "daily reminder", s.executeReminder)
	}
	if s.enabled(NotifMalam) {
		go s.runDailyJob(ctx, 21, 0, "night message", s.executeNightMessage)
	}
}

func (s *CronService) enabled(name string) bool {
	return !s.disabled[name]
}

// runDailyJob runs fn every day at hour:minute (local time).
func (s *CronService) runDailyJob(ctx context.Context, hour, minute int, label string, fn func(context.Context)) {
	for {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if !target.After(now) {
			target = target.AddDate(0, 0, 1)
		}

		log.Printf("Next %s at %v", label, target)
		select {
		case <-ctx.Done():
			return
		case <-time.After(target.Sub(now)):
			fn(ctx)
		}
	}
}

func (s *CronService) executeReminder(ctx context.Context) {
	log.Println("Executing daily reminder...")

	for _, gid := range s.groupIDs {
		groupJID, err := types.ParseJID(gid)
		if err != nil {
			log.Printf("Invalid group ID %s: %v", gid, err)
			continue
		}

		groupInfo, err := s.client.GetGroupInfo(ctx, groupJID)
		if err != nil {
			log.Printf("Failed to get group info for %s: %v", gid, err)
			continue
		}

		today := time.Now().Truncate(24 * time.Hour)
		var unreported []string
		var jids []string

		for _, participant := range groupInfo.Participants {
			phone := participant.JID.User
			userID := s.repo.ResolveLIDToPhone(ctx, phone)

			dp, err := s.repo.GetDailyProgress(ctx, userID, gid, today)
			if err == nil && (dp == nil || dp.ReportsCount == 0) {
				unreported = append(unreported, fmt.Sprintf("@%s", phone))
				jids = append(jids, participant.JID.String())
			}
		}

		if len(unreported) > 0 {
			msg := fmt.Sprintf("Assalamu'alaikum, udah jam 18:00 nih.\nAyo yang belum laporan: %s\n\nJangan lupa baca Al-Qur'an hari ini ya!\nKetik: *Alhamdulillah 1 juz*", formatList(unreported))
			s.sendToGroup(ctx, gid, msg, jids)
		}
	}
}

func (s *CronService) sendToGroup(ctx context.Context, gid string, text string, mentions []string) {
	groupJID, err := types.ParseJID(gid)
	if err != nil {
		log.Printf("Invalid group ID %s: %v", gid, err)
		return
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text:        &text,
			ContextInfo: &waE2E.ContextInfo{},
		},
	}

	if len(mentions) > 0 {
		msg.ExtendedTextMessage.ContextInfo.MentionedJID = mentions
	}

	_, err = s.client.SendMessage(ctx, groupJID, msg)
	if err != nil {
		log.Printf("Failed to send scheduled message to %s: %v", gid, err)
	}
}

func (s *CronService) broadcast(ctx context.Context, msg string) {
	for _, gid := range s.groupIDs {
		s.sendToGroup(ctx, gid, msg, nil)
	}
}

// executeMorningDigest menggabungkan countdown Ramadhan, ayat/hadist, kajian,
// dan pengingat dhuha menjadi satu pesan pagi.
func (s *CronService) executeMorningDigest(ctx context.Context) {
	log.Println("Executing morning digest...")
	s.broadcast(ctx, s.buildMorningDigest(time.Now()))
}

func (s *CronService) buildMorningDigest(now time.Time) string {
	var b strings.Builder
	b.WriteString("🌅 *Assalamu'alaikum, Shohibul Qur'an!*\n")

	if line := ramadhanLine(ramadhanStartDate(), now); line != "" {
		b.WriteString("\n" + line + "\n")
	}

	b.WriteString("\n💎 *Mutiara Hari Ini*\n")
	b.WriteString(s.motEngine.GetRandomMotivation() + "\n")

	if item := s.kajianEngine.GetRandomItem(); item != nil {
		fmt.Fprintf(&b, "\n📚 *Kajian: %s*\n%s\n", item.Title, item.Excerpt)
		if item.Link != "" {
			fmt.Fprintf(&b, "🔗 %s\n", item.Link)
		}
	}

	b.WriteString("\n☀️ Jangan lupa *Shalat Dhuha* ya")
	if keutamaan := s.prayerEngine.GetRandomKeutamaan(s.prayerEngine.GetNotificationByName("Shalat Dhuha")); keutamaan != "" {
		b.WriteString(" — " + keutamaan)
	}
	b.WriteString("\n\n📖 Selamat tilawah, semoga istiqamah hari ini! 🤲")
	return b.String()
}

// executeNightMessage mengirim konten malam Jum'at pada hari Kamis, dan
// sunnah sebelum tidur pada hari lainnya.
func (s *CronService) executeNightMessage(ctx context.Context) {
	log.Println("Executing night message...")
	s.broadcast(ctx, s.buildNightMessage(time.Now()))
}

func (s *CronService) buildNightMessage(now time.Time) string {
	if now.Weekday() == time.Thursday {
		if day := s.sunnahEngine.GetSpecialDay("Malam Jum'at"); day != nil {
			msg := fmt.Sprintf("🌙 *%s* 🌙\n\n%s\n\n", day.Name, day.Message)
			for i, item := range day.Items {
				msg += fmt.Sprintf("%d. %s\n\n", i+1, item)
			}
			msg += "Semoga kita semua dapat mengamalkannya. 🤲"
			return msg
		}
	}

	sunnah := s.sunnahEngine.GetRandomSleepSunnah()
	return fmt.Sprintf("🌙 *Sunnah Sebelum Tidur* 🌙\n\n*%s*\n\n%s\n\nSemoga kita semua istiqamah dalam mengamalkan sunnah Nabi SAW. 🤲", sunnah.Title, sunnah.Message)
}

func ramadhanStartDate() time.Time {
	if v := strings.TrimSpace(os.Getenv("RAMADHAN_DATE")); v != "" {
		if t, err := time.ParseInLocation("2006-01-02", v, time.Local); err == nil {
			return t
		}
		log.Printf("Invalid RAMADHAN_DATE %q, using default", v)
	}
	// Default: 1 Ramadan 1448 H ≈ 8 Februari 2027
	return time.Date(2027, 2, 8, 0, 0, 0, 0, time.Local)
}

// ramadhanLine returns the countdown / "hari ke-N" line for the morning digest,
// or "" when Ramadhan is more than 100 days away or already over.
func ramadhanLine(start, now time.Time) string {
	startDay := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, now.Location())
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// Round to absorb DST/zone offsets so the result is a whole day count.
	days := int(math.Round(startDay.Sub(today).Hours() / 24))

	switch {
	case days > 100:
		return ""
	case days > 1:
		return fmt.Sprintf("🌙 Insya Allah Ramadhan tinggal *%d hari lagi*. Yuk persiapkan diri!", days)
	case days == 1:
		return "🌙 Insya Allah *besok* hari pertama Ramadhan. Marhaban ya Ramadhan!"
	case days == 0:
		return "🌙 *Ramadhan Mubarak!* Hari ini hari pertama Ramadhan."
	case -days+1 <= 30:
		return fmt.Sprintf("🌙 *Ramadhan hari ke-%d* — semoga puasa dan tilawah kita diterima.", -days+1)
	default:
		return ""
	}
}

func formatList(items []string) string {
	resp := ""
	for _, item := range items {
		resp += item + " "
	}
	return resp
}
