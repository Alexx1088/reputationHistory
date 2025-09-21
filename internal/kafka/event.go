package kafka

import "time"

type ReputationEntryEvent struct {
	EventID    string    `json:"event_id"`
	UserID     string    `json:"user_id"`
	Delta      int32     `json:"delta"`
	Reason     string    `json:"reason"`
	Source     string    `json:"source"`
	OccurredAt time.Time `json:"occurred_at"`
	TraceID    string    `json:"trace_id,omitempty"`
}
