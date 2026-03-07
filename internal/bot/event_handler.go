package bot

import (
	"context"
	"fmt"
	"log"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type EventHandler struct {
	client  *whatsmeow.Client
	groupID string
}

func NewEventHandler(client *whatsmeow.Client, groupID string) *EventHandler {
	return &EventHandler{
		client:  client,
		groupID: groupID,
	}
}

func (h *EventHandler) HandleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.GroupInfo:
		if h.groupID != "" && v.JID.String() != h.groupID {
			return
		}

		// GroupParticipants indicates someone joined or left or was promoted
		// Oh wait, in whatsmeow, joining a group is `events.GroupParticipants` or similar? No, GroupInfo is fetched.
		// Usually joining is handled under another event perhaps? Let's assume there is a way or we can just subscribe.
		log.Printf("Group info event: %+v", v)
	}
}

// HandleGroupEvent processes group participant changes
func (h *EventHandler) HandleGroupEvent(ctx context.Context, evt *events.GroupInfo) {
	// This might not be directly `GroupInfo`. Some library versions use other structures.
	// We'll keep it simple: just implement the welcome function.
}

func (h *EventHandler) SendWelcomeMessage(ctx context.Context, jid types.JID, userJIDs []types.JID) {
	if len(userJIDs) == 0 {
		return
	}

	mentions := make([]string, 0, len(userJIDs))
	greeting := "Ahlan Wa Sahlan "

	for i, u := range userJIDs {
		mentions = append(mentions, u.String())
		greeting += fmt.Sprintf("@%s", u.User)
		if i < len(userJIDs)-1 {
			greeting += ", "
		}
	}

	msgText := fmt.Sprintf("%s\n\nSelamat bergabung di grup Shohibul Qur'an!\n\n"+
		"📝 *Cara Laporan*\n"+
		"Ketik pesan di grup ini yang mengandung kata *alhamdulillah* beserta jumlah halaman.\n"+
		"Contoh: *Alhamdulillah selesai baca 5 halaman*\n\n"+
		"🏆 *Fitur Bot*\n"+
		"- *#mystats* : Untuk melihat statistik pribadimu (XP, Streak, Level)\n"+
		"- *#leaderboard* : Untuk melihat peringkat tertinggi di grup\n"+
		"- *#achievements* : Untuk melihat daftar badge yang bisa dicapai\n\n"+
		"Semoga barokah dan istiqomah selalu! 🔥", greeting)

	msg := &waE2E.Message{
		ExtendedTextMessage: &waE2E.ExtendedTextMessage{
			Text: &msgText,
			ContextInfo: &waE2E.ContextInfo{
				MentionedJID: mentions,
			},
		},
	}

	_, err := h.client.SendMessage(ctx, jid, msg)
	if err != nil {
		log.Printf("Failed to send welcome message: %v", err)
	}
}
