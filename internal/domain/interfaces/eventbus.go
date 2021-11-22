package interfaces

import (
	"context"

	"github.com/leksss/banner_rotator/internal/domain/entities"
)

type EventBus interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	AddEvent(ctx context.Context, stat entities.EventStat) error
}
