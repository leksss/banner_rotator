package services

import (
	"fmt"
	"math"

	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
)

// CalculateBestBanner calculate Upper Confidence Bound algorithm
func CalculateBestBanner(log logger.Log, bannerIDs []entities.BannerID, counters entities.BannerCounterMap) entities.BannerID {
	totalShowCnt := float64(0)
	for _, id := range bannerIDs {
		if c, ok := counters[id]; ok {
			totalShowCnt += c.ShowCnt
		}
	}
	if totalShowCnt == 0 {
		totalShowCnt = float64(len(bannerIDs))
	}

	rates := make(map[entities.BannerID]float64)
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

	bestBannerID := entities.BannerID(0)
	maxRate := float64(0)
	for bannerID, rate := range rates {
		if maxRate <= rate {
			bestBannerID = bannerID
			maxRate = rate
		}
	}

	return bestBannerID
}
