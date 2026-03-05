package subscriptions

import "time"

type SubscriptionIn struct {
	UUID   string
	URL    string
	Secret string `json:"-"`
	Events []string
}

type Subscription struct {
	UUID      string
	URL       string
	Secret    string `json:"-"`
	Events    []string
	CreatedAt time.Time
	UpdatedAt time.Time
}
