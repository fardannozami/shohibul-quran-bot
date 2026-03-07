package domain

import (
	"context"
	"time"
)

type DailyProgress struct {
	UserID       string    `json:"user_id" db:"user_id"`
	GroupID      string    `json:"group_id" db:"group_id"`
	Date         time.Time `json:"date" db:"date"`
	Pages        int       `json:"pages" db:"pages"`
	ReportsCount int       `json:"reports_count" db:"reports_count"`
}

type DailyProgressRepository interface {
	GetDailyProgress(ctx context.Context, userID string, groupID string, date time.Time) (*DailyProgress, error)
	UpsertDailyProgress(ctx context.Context, progress *DailyProgress) error
	GetTotalPagesInRange(ctx context.Context, start, end time.Time, groupID string) (int, error)
}
