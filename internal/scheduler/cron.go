package scheduler

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/app/motivation"
	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
)

type CronService struct {
	client     *whatsmeow.Client
	repo       domain.BotRepository
	motEngine  *motivation.Engine
	groupID    string
}

func NewCronService(client *whatsmeow.Client, repo domain.BotRepository, motEngine *motivation.Engine, groupID string) *CronService {
	return &CronService{
		client:    client,
		repo:      repo,
		motEngine: motEngine,
		groupID:   groupID,
	}
}

func (s *CronService) Start(ctx context.Context) {
	if s.groupID == "" {
		log.Println("No group ID configured, skipping cron jobs")
		return
	}

	go s.runReminderJob(ctx)
	go s.runMotivationJob(ctx)
}

func (s *CronService) runReminderJob(ctx context.Context) {
	for {
		now := time.Now()
		target := time.Date(now.Year(), now.Month(), now.Day(), 18, 0, 0, 0, now.Location())
		if now.After(target) {
			target = target.AddDate(0, 0, 1) // next day 18:00
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
	users, err := s.repo.GetAllUsers(ctx)
	if err != nil {
		log.Printf("Failed to get users for reminder: %v", err)
		return
	}

	today := time.Now().Truncate(24 * time.Hour)
	var unreported []string
	var jids []string

	for _, u := range users {
		dp, err := s.repo.GetDailyProgress(ctx, u.ID, today)
		if err == nil && (dp == nil || dp.ReportsCount == 0) {
			unreported = append(unreported, fmt.Sprintf("@%s", u.Phone))
			jids = append(jids, u.ID)
		}
	}

	if len(unreported) > 0 {
		msg := fmt.Sprintf("Assalamu'alaikum, udah jam 18:00 nih.\nAyo yang belum lapor: %s\n\nJangan lupa baca Al-Qur'an hari ini ya!\nKetik: *Alhamdulillah 1 juz*", formatList(unreported))
		s.sendMessage(ctx, msg, jids)
	}
}

func (s *CronService) runMotivationJob(ctx context.Context) {
	for {
		now := time.Now()
		
		// generate random hour between 6 and 21
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randHour := r.Intn(16) + 6 // 6 to 21
		randMin := r.Intn(60)

		target := time.Date(now.Year(), now.Month(), now.Day(), randHour, randMin, 0, 0, now.Location())
		
		// If target is already past for today, schedule it for tomorrow
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
			s.sendMessage(ctx, msg, nil)
			
			// sleep extra an hour just to be safe it doesn't loop quickly
			time.Sleep(1 * time.Hour)
		}
	}
}

func (s *CronService) sendMessage(ctx context.Context, text string, mentions []string) {
	groupJID, err := types.ParseJID(s.groupID)
	if err != nil {
		log.Printf("Invalid group ID: %v", err)
		return
	}

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &text,
			ContextInfo: &waE2E.ContextInfo{},
		},
	}

	if len(mentions) > 0 {
		msg.ExtendedTextMessage.ContextInfo.MentionedJID = mentions
	}

	_, err = s.client.SendMessage(ctx, groupJID, msg)
	if err != nil {
		log.Printf("Failed to send scheduled message: %v", err)
	}
}

func formatList(items []string) string {
	resp := ""
	for _, item := range items {
		resp += item + " "
	}
	return resp
}
