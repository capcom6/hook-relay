package events

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

func NewEvent(eventType string, payload json.RawMessage) *Event {
	return &Event{
		Type:    eventType,
		Payload: payload,
	}
}

func (e *Event) Marshal() ([]byte, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}
	return data, nil
}

func (e *Event) Unmarshal(data []byte) error {
	err := json.Unmarshal(data, e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}
	return nil
}

type EventWrapper struct {
	Event Event
	Ack   func() bool
	Nack  func() bool
}
