package scheduler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/app/kajian"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/motivation"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/prayer"
	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

type CronService struct {
	client      *whatsmeow.Client
	repo        domain.BotRepository
	motEngine   *motivation.Engine
	prayerEngine *prayer.Engine
	kajianEngine *kajian.Engine
	groupIDs    []string
}

func NewCronService(client *whatsmeow.Client, repo domain.BotRepository, motEngine *motivation.Engine, prayerEngine *prayer.Engine, kajianEngine *kajian.Engine, groupIDs []string) *CronService {
	return &CronService{
		client:      client,
		repo:        repo,
		motEngine:   motEngine,
		prayerEngine: prayerEngine,
		kajianEngine: kajianEngine,
		groupIDs:    groupIDs,
	}
}

func (s *CronService) Start(ctx context.Context) {
	if len(s.groupIDs) == 0 {
		log.Println("No group IDs configured, skipping cron jobs")
		return
	}

	go s.runReminderJob(ctx)
	go s.runMotivationJob(ctx)
	go s.runPrayerNotificationJob(ctx)
	go s.runKajianJob(ctx)
}

func (s *CronService) runReminderJob(ctx context.Context) {
	for {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
		if now.After(target) {
			target = target.AddDate(0, 0, 1)
		}

		duration := target.Sub(now)
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			s.executeReminder(ctx)
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

func (s *CronService) runMotivationJob(ctx context.Context) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		now := time.Now()

		randHour := r.Intn(16) + 6
		randMin := r.Intn(60)

		target := time.Date(now.Year(), now.Month(), now.Day(), randHour, randMin, 0, 0, now.Location())

		if target.Before(now) {
			target = target.AddDate(0, 0, 1)
		}

		duration := target.Sub(now)
		log.Printf("Next motivation will be sent at %v", target)

		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			quote := s.motEngine.GetRandomMotivation()
			msg := fmt.Sprintf("🌟 *Daily Motivation* 🌟\n\n%s", quote)
			for _, gid := range s.groupIDs {
				s.sendToGroup(ctx, gid, msg, nil)
			}

			now2 := time.Now()
			nextMidnight := time.Date(now2.Year(), now2.Month(), now2.Day()+1, 0, 0, 0, 0, now2.Location())

			select {
			case <-ctx.Done():
				return
			case <-time.After(nextMidnight.Sub(now2)):
			}
		}
	}
}

func (s *CronService) runPrayerNotificationJob(ctx context.Context) {
	for {
		now := time.Now()

		// Find next prayer notification time
		nextPrayerTime := s.findNextPrayerTime(now)
		duration := nextPrayerTime.Sub(now)

		log.Printf("Next prayer notification at %v", nextPrayerTime)

		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			s.executePrayerNotification(ctx, nextPrayerTime)
		}
	}
}

func (s *CronService) findNextPrayerTime(now time.Time) time.Time {
	notifications := s.prayerEngine.GetAllNotifications()

	// Parse all notification times and find the next one
	for _, n := range notifications {
		hour, minute := parseTime(n.Time)
		target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, 0, 0, now.Location())
		if target.After(now) {
			return target
		}
	}

	// If all times have passed today, schedule for tomorrow's first notification
	firstNotification := notifications[0]
	hour, minute := parseTime(firstNotification.Time)
	target := time.Date(now.Year(), now.Month(), now.Day()+1, hour, minute, 0, 0, now.Location())
	return target
}

func (s *CronService) executePrayerNotification(ctx context.Context, scheduledTime time.Time) {
	hour := scheduledTime.Hour()
	minute := scheduledTime.Minute()

	notification := s.prayerEngine.GetNotificationByTime(hour, minute)
	if notification == nil {
		return
	}

	log.Printf("Executing prayer notification: %s", notification.Name)

	keutamaan := s.prayerEngine.GetRandomKeutamaan(notification)
	msg := fmt.Sprintf("%s\n\n*Keutamaan %s:*\n%s\n\nSemoga kita semua dilimpahkan keberkahan dan istiqamah dalam beribadah. 🤲", notification.Message, notification.Name, keutamaan)

	for _, gid := range s.groupIDs {
		s.sendToGroup(ctx, gid, msg, nil)
	}
}

func parseTime(timeStr string) (int, int) {
	var hour, minute int
	fmt.Sscanf(timeStr, "%d:%d", &hour, &minute)
	return hour, minute
}

func (s *CronService) runKajianJob(ctx context.Context) {
	for {
		now := time.Now()
		// Schedule kajian for 07:00 every day
		target := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, now.Location())
		if now.After(target) {
			target = target.AddDate(0, 0, 1)
		}

		duration := target.Sub(now)
		log.Printf("Next kajian notification at %v", target)

		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			s.executeKajianNotification(ctx)
		}
	}
}

func (s *CronService) executeKajianNotification(ctx context.Context) {
	log.Println("Executing daily kajian notification...")

	msg := s.kajianEngine.GetRandomKajian()
	for _, gid := range s.groupIDs {
		s.sendToGroup(ctx, gid, msg, nil)
	}
}

func (s *CronService) sendMessage(ctx context.Context, text string, mentions []string) {
	for _, gid := range s.groupIDs {
		groupJID, err := types.ParseJID(gid)
		if err != nil {
			log.Printf("Invalid group ID %s: %v", gid, err)
			continue
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
}

func formatList(items []string) string {
	resp := ""
	for _, item := range items {
		resp += item + " "
	}
	return resp
}
