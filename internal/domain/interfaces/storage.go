package interfaces

import (
	"context"

	"github.com/leksss/banner_rotator/internal/domain/entities"
)

type Storage interface {
	AddBanner(ctx context.Context, slotID, bannerID uint64) error
	RemoveBanner(ctx context.Context, slotID, bannerID uint64) error
	IncrementHit(ctx context.Context, slotID, bannerID, groupID uint64) error
	IncrementShow(ctx context.Context, slotID, bannerID, groupID uint64) error
	GetBannersBySlot(ctx context.Context, slotID uint64) ([]entities.BannerID, error)
	GetSlotCounters(ctx context.Context, slotID, groupID uint64) (entities.BannerCounterMap, error)
}
