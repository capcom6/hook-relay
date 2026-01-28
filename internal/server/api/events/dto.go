package events

import "encoding/json"

type EventRequest struct {
	ID      string          `json:"id"      validate:"required,max=64"`
	Type    string          `json:"type"    validate:"required,max=64"`
	Payload json.RawMessage `json:"payload" validate:"max=65536"`
}
