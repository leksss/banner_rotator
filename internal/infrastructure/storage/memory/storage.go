package memorystorage

import (
	"context"
	"sync"
)

type Storage struct {
	mu      sync.RWMutex //nolint
	banners []uint64
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return nil
}

func (s *Storage) AddBanner(ctx context.Context, slotID, bannerID uint64) error {
	s.banners = append(s.banners, bannerID)
	return nil
}

func (s *Storage) RemoveBanner(ctx context.Context, slotID, bannerID uint64) error {
	for i, banID := range s.banners {
		if banID == bannerID {
			s.banners = append(s.banners[:i], s.banners[i+1:]...)
			break
		}
	}
	return nil
}

func (s *Storage) HitBanner(ctx context.Context, slotID, bannerID, groupID uint64) error {
	return nil
}

func (s *Storage) GetBanner(ctx context.Context, slotID, groupID uint64) (uint64, error) {
	if len(s.banners) == 0 {
		return 0, nil
	}
	return s.banners[0], nil
}
