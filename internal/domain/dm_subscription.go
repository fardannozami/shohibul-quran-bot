package domain

import "context"

// DMSubscriptionRepository stores users who opted in to receive bot messages
// via personal chat. Opt-in only happens when the user DMs the bot first, so
// the bot never initiates a conversation with someone who didn't ask for it.
type DMSubscriptionRepository interface {
	SetDMSubscription(ctx context.Context, userID, chatJID string) error
	DeleteDMSubscription(ctx context.Context, userID string) error
	// GetDMSubscription returns the chat JID to DM, or "" if not subscribed.
	GetDMSubscription(ctx context.Context, userID string) (string, error)
}
