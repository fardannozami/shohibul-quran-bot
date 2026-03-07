package domain

import (
	"context"
	"time"
)

type ReportLog struct {
	ID      string    `json:"id" db:"id"`
	UserID  string    `json:"user_id" db:"user_id"`
	Pages   int       `json:"pages" db:"pages"`
	Message string    `json:"message" db:"message"`
	Date    time.Time `json:"date" db:"date"`
}

type ReportLogRepository interface {
	InsertReport(ctx context.Context, report *ReportLog) error
	GetReportsByUser(ctx context.Context, userID string) ([]*ReportLog, error)
}
