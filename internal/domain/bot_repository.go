package domain

import "context"

type BotRepository interface {
	UserRepository
	ReportLogRepository
	DailyProgressRepository
	BadgeRepository
	InitDatabase(ctx context.Context) error
}
