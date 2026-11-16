package domain

import "context"

type BotRepository interface {
	UserRepository
	ReportLogRepository
	DailyProgressRepository
	BadgeRepository
	DMSubscriptionRepository
	InitDatabase(ctx context.Context) error
}
