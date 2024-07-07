package repository

import (
	"context"
	"time"

	"github.com/mymmrac/telego"
)

type ScheduledMessage struct {
	Message     string
	TriggerTime time.Time
}

type MessageDelivery struct {
	ScheduledMessage ScheduledMessage
	SentMessages     map[int64]*telego.Message
	// sets to true when all subscribed users are listed in SentMessages. After that no triggers will be called
	FullyDelivered bool
}

type RepositoryUser interface {
	Remember(ctx context.Context, user telego.User) error
	GetRememberedUser(ctx context.Context, userID int64) (telego.User, error)
	ListRememberedUser(ctx context.Context) ([]telego.User, error)

	Subscribe(ctx context.Context, userID int64) error
	Unsubscribe(ctx context.Context, userID int64) error

	ScheduleMessage(ctx context.Context, messageSchedule ScheduledMessage) error
	ListScheduledMessages(ctx context.Context) ([]MessageDelivery, error)
}
