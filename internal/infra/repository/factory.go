package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/fardannozami/shohibul-quran-bot/internal/config"
	"github.com/fardannozami/shohibul-quran-bot/internal/domain"
	"github.com/fardannozami/shohibul-quran-bot/internal/infra/sqlite"
	_ "modernc.org/sqlite"
)

func NewBotRepository(cfg config.Config) domain.BotRepository {
	log.Println("Using SQLite database")
	// Enable WAL mode and busy timeout to avoid "database is locked" errors
	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)", cfg.SQLitePath)
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	repo := sqlite.NewBotRepository(db)
	// Initialize table if needed
	if err := repo.InitDatabase(context.Background()); err != nil {
		log.Printf("Failed to init database: %v", err)
	}

	return repo
}
