package interfaces

import (
	"context"

	"github.com/leksss/banner_rotator/internal/domain/entities"
)

type DatabaseConf struct {
	Host     string
	User     string
	Password string
	Name     string
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	AddBanner(ctx context.Context, slotID, bannerID uint64) error
	RemoveBanner(ctx context.Context, slotID, bannerID uint64) error
	IncrementHit(ctx context.Context, slotID, bannerID, groupID uint64) error
	IncrementShow(ctx context.Context, slotID, bannerID, groupID uint64) error
	GetBannersBySlot(ctx context.Context, slotID uint64) ([]uint64, error)
	GetSlotCounters(ctx context.Context, slotID, groupID uint64) (map[uint64]*entities.Counter, error)
}
