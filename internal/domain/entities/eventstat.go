package entities

import "time"

const (
	EventTypeUndefined = 0
	EventTypeHit       = 1
	EventTypeShow      = 2
)

type EventStat struct {
	EventType uint64
	SlotID    uint64
	BannerID  uint64
	GroupID   uint64
	CreatedAt time.Time
}
