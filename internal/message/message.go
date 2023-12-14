package message

import (
	"github.com/google/uuid"
)

type Message struct {
	DeviceID *uuid.UUID `json:"device_id,omitempty"`
	ID       uuid.UUID  `json:"id"`
	Kind     int        `json:"kind"`
	Message  string     `json:"message"`
}
