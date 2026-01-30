package subscriptions

import (
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/bun"
)

// subscriptionModel represents a webhook subscription in the database.
type subscriptionModel struct {
	bun.BaseModel `bun:"table:subscriptions"`

	ID        int64     `bun:"id,pk,nullzero"`
	UUID      string    `bun:"uuid"`
	URL       string    `bun:"url"`
	Secret    string    `bun:"secret,nullzero"`
	CreatedAt time.Time `bun:"created_at,nullzero"`
	UpdatedAt time.Time `bun:"updated_at,nullzero"`

	Events []eventModel `bun:"rel:has-many,join:id=subscription_id"`
}

func newSubscriptionModel(subscription SubscriptionIn) *subscriptionModel {
	return &subscriptionModel{
		BaseModel: bun.BaseModel{},
		ID:        0,
		UUID:      subscription.UUID,
		URL:       subscription.URL,
		Secret:    subscription.Secret,
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		Events:    lo.Map(subscription.Events, func(item string, _ int) eventModel { return *newEventModel(item) }),
	}
}

func (s *subscriptionModel) toDomain() *Subscription {
	if s == nil {
		return nil
	}

	return &Subscription{
		UUID:      s.UUID,
		URL:       s.URL,
		Secret:    s.Secret,
		Events:    lo.Map(s.Events, func(item eventModel, _ int) string { return item.Event }),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

// eventModel represents a webhook event in the database.
type eventModel struct {
	bun.BaseModel `bun:"table:subscription_events"`

	ID             int64  `bun:"id,pk,nullzero"`
	SubscriptionID int64  `bun:"subscription_id"`
	Event          string `bun:"event"`
}

func newEventModel(event string) *eventModel {
	return &eventModel{
		BaseModel:      bun.BaseModel{},
		Event:          event,
		ID:             0,
		SubscriptionID: 0,
	}
}
