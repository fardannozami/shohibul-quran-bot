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
	client    *whatsmeow.Client
	repo      domain.BotRepository
	motEngine *motivation.Engine
	groupIDs  []string
}

func NewCronService(client *whatsmeow.Client, repo domain.BotRepository, motEngine *motivation.Engine, groupIDs []string) *CronService {
	return &CronService{
		client:    client,
		repo:      repo,
		motEngine: motEngine,
		groupIDs:  groupIDs,
	}
}

func (s *CronService) Start(ctx context.Context) {
	if len(s.groupIDs) == 0 {
		log.Println("No group IDs configured, skipping cron jobs")
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
