package interfaces

import (
	"context"

	"github.com/leksss/banner_rotator/internal/domain/entities"
)

type EventBus interface {
	AddEvent(ctx context.Context, stat entities.EventStat) error
}
