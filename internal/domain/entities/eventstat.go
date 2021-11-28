package entities

import "time"

type EventType int

const (
	EventTypeUndefined EventType = iota
	EventTypeHit
	EventTypeShow
)

type EventStat struct {
	EventType EventType
	SlotID    uint64
	BannerID  uint64
	GroupID   uint64
	CreatedAt time.Time
}
