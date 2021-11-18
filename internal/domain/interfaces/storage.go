package interfaces

import (
	"context"
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
	HitBanner(ctx context.Context, slotID, bannerID, groupID uint64) error
	GetBanner(ctx context.Context, slotID, groupID uint64) (uint64, error)
}
