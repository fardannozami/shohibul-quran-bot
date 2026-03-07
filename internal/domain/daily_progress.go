package domain

import (
	"context"
	"time"
)

type DailyProgress struct {
	UserID       string    `json:"user_id" db:"user_id"`
	Date         time.Time `json:"date" db:"date"`
	Pages        int       `json:"pages" db:"pages"`
	ReportsCount int       `json:"reports_count" db:"reports_count"`
}

type DailyProgressRepository interface {
	GetDailyProgress(ctx context.Context, userID string, date time.Time) (*DailyProgress, error)
	UpsertDailyProgress(ctx context.Context, progress *DailyProgress) error
}
