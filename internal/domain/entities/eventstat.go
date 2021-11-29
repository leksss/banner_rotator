package entities

import "time"

type EventType int

const (
	// EventTypeHit hit event.
	EventTypeHit EventType = iota + 1
	// EventTypeShow show event.
	EventTypeShow
)

type EventStat struct {
	EventType EventType
	SlotID    uint64
	BannerID  uint64
	GroupID   uint64
	CreatedAt time.Time
}
