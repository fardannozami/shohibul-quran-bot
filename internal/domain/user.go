package domain

import (
	"context"
	"time"
)

type User struct {
	ID       string    `json:"id" db:"id"`
	GroupID  string    `json:"group_id" db:"group_id"`
	Phone    string    `json:"phone" db:"phone"`
	Name     string    `json:"name" db:"name"`
	XP       int       `json:"xp" db:"xp"`
	Level    int       `json:"level" db:"level"`
	Streak   int       `json:"streak" db:"streak"`
	JoinedAt    time.Time `json:"joined_at" db:"joined_at"`
	DailyTarget int       `json:"daily_target" db:"daily_target"`
}

type UserRepository interface {
	GetUser(ctx context.Context, id string, groupID string) (*User, error)
	CreateUser(ctx context.Context, user *User) error
	UpdateUser(ctx context.Context, user *User) error
	GetAllUsers(ctx context.Context, groupID string) ([]*User, error)
	ResolveLIDToPhone(ctx context.Context, lid string) string
}
