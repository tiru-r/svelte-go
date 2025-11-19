package types

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event represents a domain event in the system
type Event struct {
	ID          string                 `json:"id"`
	Type        string                 `json:"type"`
	Source      string                 `json:"source"`
	AggregateID string                 `json:"aggregate_id,omitempty"`
	Version     int                    `json:"version"`
	Data        map[string]interface{} `json:"data"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
}

// NewEvent creates a new event with generated ID and timestamp
func NewEvent(eventType, source string, data map[string]interface{}) *Event {
	return &Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    source,
		Version:   1,
		Data:      data,
		Metadata:  make(map[string]interface{}),
		Timestamp: time.Now(),
	}
}

// WithAggregateID sets the aggregate ID for domain-driven design
func (e *Event) WithAggregateID(id string) *Event {
	e.AggregateID = id
	return e
}

// WithMetadata adds metadata to the event
func (e *Event) WithMetadata(key string, value interface{}) *Event {
	if e.Metadata == nil {
		e.Metadata = make(map[string]interface{})
	}
	e.Metadata[key] = value
	return e
}

// ToJSON marshals the event to JSON
func (e *Event) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// EventHandler is a function that processes events
type EventHandler func(*Event) error

// EventBus interface for event publishing and subscribing
type EventBus interface {
	Publish(subject string, event *Event) error
	Subscribe(subject string, handler EventHandler) error
	SubscribeQueue(subject, queue string, handler EventHandler) error
	Close() error
}