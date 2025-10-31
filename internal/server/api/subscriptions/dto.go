package subscriptions

import (
	"github.com/capcom6/hook-relay/internal/subscriptions"
	"github.com/capcom6/hook-relay/pkg/privacy"
)

type Subscription struct {
	UUID   string               `json:"uuid"             validate:"required,uuid"`
	URL    string               `json:"url"              validate:"required,https_url,max=256"`
	Secret privacy.HiddenString `json:"secret,omitempty" validate:"omitempty,max=256"`
	Events []string             `json:"events"           validate:"required,min=1,dive,required,max=32"`
}

func newSubscription(source *subscriptions.Subscription) *Subscription {
	return &Subscription{
		UUID:   source.UUID,
		URL:    source.URL,
		Secret: "",
		Events: source.Events,
	}
}
