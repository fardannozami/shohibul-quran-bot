package bot

import (
	"context"
	"fmt"
	"log"

	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
)

type EventHandler struct {
	client   *whatsmeow.Client
	groupIDs []string

	mu             sync.Mutex
	pendingWelcome map[types.JID][]types.JID
	welcomeTimer   *time.Timer
}

func NewEventHandler(client *whatsmeow.Client, groupIDs []string) *EventHandler {
	return &EventHandler{
		client:         client,
		groupIDs:        groupIDs,
		pendingWelcome: make(map[types.JID][]types.JID),
	}
}

func (h *EventHandler) HandleEvent(evt interface{}) {
	switch v := evt.(type) {
	case *events.GroupInfo:
		if len(h.groupIDs) > 0 {
			allowed := false
			for _, gid := range h.groupIDs {
				if v.JID.String() == gid {
					allowed = true
					break
				}
			}
			if !allowed {
				return
			}
		}

		// GroupParticipants indicates someone joined or left or was promoted
		// Oh wait, in whatsmeow, joining a group is `events.GroupParticipants` or similar? No, GroupInfo is fetched.
		// Usually joining is handled under another event perhaps? Let's assume there is a way or we can just subscribe.
		log.Printf("Group info event: %+v", v)
	}
}

func (h *EventHandler) QueueWelcomeMessage(ctx context.Context, groupJID types.JID, userJIDs []types.JID) {
	if len(userJIDs) == 0 {
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Append newbies to the specific group
	for _, newU := range userJIDs {
		exists := false
		for _, existingU := range h.pendingWelcome[groupJID] {
			if existingU.String() == newU.String() {
				exists = true
				break
			}
		}
		if !exists {
			h.pendingWelcome[groupJID] = append(h.pendingWelcome[groupJID], newU)
		}
	}

	// Reset the timer
	if h.welcomeTimer != nil {
		h.welcomeTimer.Stop()
	}

	// Wait 10 seconds for more people to join before sending
	h.welcomeTimer = time.AfterFunc(10*time.Second, func() {
		h.mu.Lock()
		// Copy the pending welcomes to process them
		processMap := make(map[types.JID][]types.JID)
		for gJID, users := range h.pendingWelcome {
			processMap[gJID] = users
		}
		// Clear pending
		h.pendingWelcome = make(map[types.JID][]types.JID)
		h.welcomeTimer = nil
		h.mu.Unlock()

		// Send messages outside the lock
		for gJID, users := range processMap {
			if len(users) > 0 {
				h.SendWelcomeMessage(context.Background(), gJID, users)
			}
		}
	})
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

	msgText := greeting + "\n\n" +
		"Assalamu'alaikum warahmatullahi wabarakatuh 👋\n\n" +
		"Selamat datang di grup *Shohibul Qur'an*\n\n" +
		"Semoga kita semua dimudahkan untuk istiqomah membaca Al-Qur'an setiap hari.\n\n" +
		"📌 *Cara laporan membaca Qur'an:*\n\n" +
		"Cukup kirim pesan dengan kata:\n" +
		"\"Alhamdulillah\"\n\n" +
		"Contoh:\n" +
		"Alhamdulillah sudah baca 2 halaman\n" +
		"Alhamdulillah 1 juz\n" +
		"Alhamdulillah sudah mengaji hari ini\n\n" +
		"Alhamdulillah sudah baca al baqoroh 1 - 30\n\n" +
		"Bot akan otomatis mencatat laporan Anda.\n\n" +
		"⏰ Jika sampai jam 18:00 belum laporan, bot akan mengingatkan.\n\n" +
		"🎯 *Target kita:*\n" +
		"Istiqomah membaca Al-Qur'an setiap hari walaupun hanya beberapa ayat.\n\n" +
		"Rasulullah ﷺ bersabda:\n\n" +
		"\"Sebaik-baik kalian adalah yang mempelajari Al-Qur'an dan mengajarkannya.\"\n" +
		"(HR. Bukhari)\n\n" +
		"Barakallahu fiikum 🤍\n\n"

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
