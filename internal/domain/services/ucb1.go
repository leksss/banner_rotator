package services

import (
	"fmt"
	"math"

	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
)

// calculate Upper Confidence Bound algorithm
func CalculateBestBanner(log logger.Log, bannerIDs []uint64, counters map[uint64]*entities.Counter) uint64 {
	totalShowCnt := float64(0)
	for _, id := range bannerIDs {
		if c, ok := counters[id]; ok {
			totalShowCnt += c.ShowCnt
		}
	}
	if totalShowCnt == 0 {
		totalShowCnt = float64(len(bannerIDs))
	}

	rates := make(map[uint64]float64)
	for _, id := range bannerIDs {
		c, ok := counters[id]
		if !ok {
			c = &entities.Counter{
				HitCnt:  0,
				ShowCnt: 1,
			}
		}
		rates[id] = c.HitCnt/c.ShowCnt + math.Sqrt((2*math.Log(totalShowCnt))/c.ShowCnt)
		log.Info(fmt.Sprintf("rate for bannerID %d: %f", id, rates[id]))
	}

	bestBannerID := uint64(0)
	maxRate := float64(0)
	for bannerID, rate := range rates {
		if maxRate <= rate {
			bestBannerID = bannerID
			maxRate = rate
		}
	}

	return bestBannerID
}
