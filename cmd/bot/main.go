package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fardannozami/shohibul-quran-bot/internal/app/ai"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/gamification"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/kajian"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/motivation"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/prayer"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/sunnah"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/tafsir"
	"github.com/fardannozami/shohibul-quran-bot/internal/app/usecase"
	"github.com/fardannozami/shohibul-quran-bot/internal/bot"
	"github.com/fardannozami/shohibul-quran-bot/internal/config"
	"github.com/fardannozami/shohibul-quran-bot/internal/infra/repository"
	"github.com/fardannozami/shohibul-quran-bot/internal/infra/wa"
	"github.com/fardannozami/shohibul-quran-bot/internal/parser"
	"github.com/fardannozami/shohibul-quran-bot/internal/scheduler"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	walog "go.mau.fi/whatsmeow/util/log"
)

func main() {
	// 1. Load Config
	cfg := config.Load()

	// 2. Logger
	logger := walog.Stdout("Client", "INFO", true)

	// 3. Database & Repositories
	repo := repository.NewBotRepository(cfg)

	// 4. Use Cases
	parserMod := parser.NewReportParser()
	aiEngine := ai.NewEngine(cfg.GeminiAPIKey)
	tafsirEngine := tafsir.NewEngine()
	tafsirEngine.SetAI(aiEngine)
	gameEngine := gamification.NewEngine(repo)
	gameEngine.SetTafsir(tafsirEngine)
	motEngine := motivation.NewEngine()
	prayerEngine := prayer.NewEngine()
	kajianEngine := kajian.NewEngine()
	sunnahEngine := sunnah.NewEngine()
	handleMessageUC := usecase.NewHandleMessageUsecase(repo, parserMod, gameEngine, motEngine)

	// 5. WhatsApp Service
	waService := wa.NewService(cfg.SQLitePath, logger, cfg.SupabaseURL, cfg.SupabaseKey)

	// 5.5 Wait until client is ready or expose it to handlers
	welcomeHandler := bot.NewEventHandler(nil, cfg.GroupIDs) // We will inject client later after waService is initialized.
	cronService := scheduler.NewCronService(nil, repo, motEngine, prayerEngine, kajianEngine, sunnahEngine, cfg.GroupIDs)

	// 6. Register Message Handler
	waService.SetMessageHandler(func(ctx context.Context, client *whatsmeow.Client, evt *events.Message) {
		// Log all incoming messages with their Chat ID (useful for getting groupID)
		fmt.Printf("[DEBUG] Incoming message from Chat ID: %s (IsGroup: %v)\n", evt.Info.Chat.String(), evt.Info.IsGroup)
		if evt.Info.IsGroup {
			fmt.Printf("💡 [INFO] Group ID is: %s\n", evt.Info.Chat.String())
		}

		// Only handle messages from groups or specific sources if needed.
		// For now, we filter by GroupIDs if configured.
		if len(cfg.GroupIDs) > 0 {
			allowed := false
			for _, gid := range cfg.GroupIDs {
				if evt.Info.Chat.String() == gid {
					allowed = true
					break
				}
			}
			if !allowed {
				return
			}
		}

		// Ignore messages from self
		if evt.Info.IsFromMe {
			return
		}

		// Get sender info - resolve LID to phone number for consistent user tracking
		senderJID := evt.Info.Sender
		var userID string
		if senderJID.Server == "lid" || senderJID.Server == types.DefaultUserServer && len(senderJID.User) > 15 {
			// Looks like a LID, try to resolve to phone number
			userID = repo.ResolveLIDToPhone(ctx, senderJID.User)
		} else {
			// Already a phone number
			userID = senderJID.User
		}

		pushName := evt.Info.PushName
		if pushName == "" {
			pushName = "Unknown" // Fallback name
		}

		// Get message content
		msg := ""
		if evt.Message.Conversation != nil {
			msg = *evt.Message.Conversation
		} else if evt.Message.ExtendedTextMessage != nil && evt.Message.ExtendedTextMessage.Text != nil {
			msg = *evt.Message.ExtendedTextMessage.Text
		} else if evt.Message.ImageMessage != nil && evt.Message.ImageMessage.Caption != nil {
			msg = *evt.Message.ImageMessage.Caption
		} else if evt.Message.VideoMessage != nil && evt.Message.VideoMessage.Caption != nil {
			msg = *evt.Message.VideoMessage.Caption
		} else if evt.Message.DocumentMessage != nil && evt.Message.DocumentMessage.Caption != nil {
			msg = *evt.Message.DocumentMessage.Caption
		}

		if msg == "" {
			return
		}

		fmt.Printf("Message from %s (%s): %s\n", pushName, userID, msg)

		if isReplyToBot(evt, client) {
			var response string
			if cfg.GeminiAPIKey != "" {
				r, err := aiEngine.GenerateResponse(ctx, msg)
				if err != nil {
					log.Printf("AI response error: %v", err)
					response = "Maaf, terjadi kendala saat menghubungi AI. Coba lagi beberapa saat ya 🙏"
				} else if r != "" {
					response = r
				} else {
					response = "Maaf, AI belum bisa menjawab saat ini. Coba lagi ya 🙏"
				}
			} else {
				response = "Maaf, fitur AI sedang tidak aktif. Silakan hubungi admin untuk mengaktifkannya."
			}
			if response != "" {
				applyReplyDelay(cfg, waService, ctx, evt.Info.Chat)
				resp := &waE2E.Message{Conversation: &response}
				if _, err := waService.GetClient().SendMessage(ctx, evt.Info.Chat, resp); err != nil {
					log.Printf("Failed to send AI response: %v", err)
				}
			}
			return
		}

		// Execute Use Case
		response, err := handleMessageUC.Execute(ctx, userID, pushName, msg, evt.Info.Chat.String())
		if err != nil {
			log.Printf("Error handling message: %v", err)
			return
		}

		if response != "" {
			// Apply reply delay to appear more human-like
			delayMs := cfg.ReplyDelayMinMs
			if cfg.ReplyDelayMaxMs > cfg.ReplyDelayMinMs {
				// Random delay between min and max
				delayMs = cfg.ReplyDelayMinMs + rand.Intn(cfg.ReplyDelayMaxMs-cfg.ReplyDelayMinMs+1)
			}

			if delayMs > 0 {
				// Show typing indicator if enabled
				if cfg.ShowTyping {
					_ = waService.GetClient().SendChatPresence(ctx, evt.Info.Chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
				}

				log.Printf("Delaying reply by %dms", delayMs)
				time.Sleep(time.Duration(delayMs) * time.Millisecond)

				// Clear typing indicator
				if cfg.ShowTyping {
					_ = waService.GetClient().SendChatPresence(ctx, evt.Info.Chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
				}
			}

			// Send response
			resp := &waE2E.Message{
				Conversation: &response,
			}
			_, err := waService.GetClient().SendMessage(ctx, evt.Info.Chat, resp)
			if err != nil {
				log.Printf("Failed to send response: %v", err)
			}
		}
	})
	
	// 7. Initialize Client (DB, Device, etc) - DO NOT CONNECT YET
	if err := waService.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize WhatsApp service: %v", err)
	}

	waService.GetClient().AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.GroupInfo:
			// Action mapping from v.Join (list of users who joined)
			if len(v.Join) > 0 {
				if len(cfg.GroupIDs) > 0 {
					allowed := false
					for _, gid := range cfg.GroupIDs {
						if v.JID.String() == gid {
							allowed = true
							break
						}
					}
					if !allowed {
						return
					}
				}
				
				// Run in background to avoid blocking event handler
				go func() {
					welcomeHandler.QueueWelcomeMessage(context.Background(), v.JID, v.Join)
				}()
			}
		}
	})

	welcomeHandler = bot.NewEventHandler(waService.GetClient(), cfg.GroupIDs)
	cronService = scheduler.NewCronService(waService.GetClient(), repo, motEngine, prayerEngine, kajianEngine, sunnahEngine, cfg.GroupIDs)
	
	// Start Scheduler
	cronCtx, cronCancel := context.WithCancel(context.Background())
	defer cronCancel()
	cronService.Start(cronCtx)

	// 8. Connect / Login Logic
	if !waService.IsLoggedIn() {
		if cfg.BotPhone != "" {
			// Pair Code Mode
			// Must connect first to pair
			if err := waService.Connect(); err != nil {
				log.Fatalf("Failed to connect for pairing: %v", err)
			}

			log.Println("Not logged in. Attempting to pair with phone:", cfg.BotPhone)
			code, err := waService.Pair(cfg.BotPhone)
			if err != nil {
				log.Printf("Failed to generate pair code: %v", err)
			} else {
				log.Println("==================================================")
				log.Printf("PAIR CODE: %s", code)
				log.Println("==================================================")
				log.Println("Please verify this code on your WhatsApp (Linked Devices > Link with phone number)")
			}
		} else {
			// QR Code Mode
			log.Println("Not logged in. BOT_PHONE not set. Printing QR...")
			// PrintQR handles GetQRChannel AND Connect() internally to ensure no race condition
			waService.PrintQR()
		}
	} else {
		// Already logged in, just connect
		if err := waService.Connect(); err != nil {
			log.Fatalf("Failed to connect: %v", err)
		}
		log.Println("Client is already logged in.")
	}

	log.Println("Bot is running... Press Ctrl+C to exit.")

	// 8. Wait for OS Signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down...")
	waService.Disconnect()
	os.Exit(0)
}

func getContextInfo(msg *waE2E.Message) *waE2E.ContextInfo {
	if msg == nil {
		return nil
	}
	if m := msg.ExtendedTextMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.ImageMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.VideoMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.DocumentMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.AudioMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.StickerMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.ContactMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.LocationMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.ButtonsMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	if m := msg.ListMessage; m != nil && m.ContextInfo != nil {
		return m.ContextInfo
	}
	return nil
}

func isReplyToBot(evt *events.Message, client *whatsmeow.Client) bool {
	ctxInfo := getContextInfo(evt.Message)
	if ctxInfo == nil || ctxInfo.StanzaID == nil || *ctxInfo.StanzaID == "" {
		return false
	}
	if ctxInfo.Participant == nil || *ctxInfo.Participant == "" {
		return false
	}
	if client == nil || client.Store == nil {
		return false
	}
	targetStr := *ctxInfo.Participant
	targetUser := targetStr
	if idx := strings.Index(targetUser, "@"); idx > 0 {
		targetUser = targetUser[:idx]
	}
	botUsers := map[string]struct{}{}
	if client.Store.ID != nil && client.Store.ID.User != "" {
		botUsers[client.Store.ID.User] = struct{}{}
	}
	if lid := client.Store.GetLID(); lid.User != "" {
		botUsers[lid.User] = struct{}{}
	}
	if _, ok := botUsers[targetUser]; ok {
		return true
	}
	if parsed, err := types.ParseJID(targetStr); err == nil {
		if parsed.Server == types.HiddenUserServer {
			if pn, err := client.Store.LIDs.GetPNForLID(context.Background(), parsed); err == nil && pn.User != "" {
				if _, ok := botUsers[pn.User]; ok {
					return true
				}
			}
		} else if parsed.Server == types.DefaultUserServer {
			if lidJID, err := client.Store.LIDs.GetLIDForPN(context.Background(), parsed); err == nil && lidJID.User != "" {
				if _, ok := botUsers[lidJID.User]; ok {
					return true
				}
			}
		}
	}
	return false
}

// applyReplyDelay adds a configurable human-like delay (and optional typing indicator) before replying.
func applyReplyDelay(cfg config.Config, waService *wa.Service, ctx context.Context, chat types.JID) {
	delayMs := cfg.ReplyDelayMinMs
	if cfg.ReplyDelayMaxMs > cfg.ReplyDelayMinMs {
		delayMs = cfg.ReplyDelayMinMs + rand.Intn(cfg.ReplyDelayMaxMs-cfg.ReplyDelayMinMs+1)
	}

	if delayMs > 0 {
		if cfg.ShowTyping {
			_ = waService.GetClient().SendChatPresence(ctx, chat, types.ChatPresenceComposing, types.ChatPresenceMediaText)
		}

		log.Printf("Delaying reply by %dms", delayMs)
		time.Sleep(time.Duration(delayMs) * time.Millisecond)

		if cfg.ShowTyping {
			_ = waService.GetClient().SendChatPresence(ctx, chat, types.ChatPresencePaused, types.ChatPresenceMediaText)
		}
	}
}
