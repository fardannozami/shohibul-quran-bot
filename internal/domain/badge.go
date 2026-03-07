package domain

import (
	"context"
	"time"
)

type BadgeLog struct {
	UserID    string    `json:"user_id" db:"user_id"`
	GroupID   string    `json:"group_id" db:"group_id"`
	Badge     string    `json:"badge" db:"badge"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type BadgeRepository interface {
	InsertBadge(ctx context.Context, badge *BadgeLog) error
	GetBadgesByUser(ctx context.Context, userID string, groupID string) ([]*BadgeLog, error)
}
