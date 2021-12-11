package services

import (
	"math"

	"github.com/leksss/banner_rotator/internal/domain/entities"
)

// CalculateBestBanner calculate Upper Confidence Bound algorithm.
func CalculateBestBanner(bannerIDs []entities.BannerID, counters entities.BannerCounterMap) entities.BannerID {
	totalShowCnt := float64(0)
	for _, id := range bannerIDs {
		if c, ok := counters[id]; ok {
			totalShowCnt += c.ShowCnt
		}
	}
	if totalShowCnt == 0 {
		totalShowCnt = float64(len(bannerIDs))
	}

	maxRate := float64(0)
	bestBannerID := entities.BannerID(0)
	for _, bannerID := range bannerIDs {
		c, ok := counters[bannerID]
		if !ok || c.ShowCnt == 0 {
			c = entities.Counter{HitCnt: 0, ShowCnt: 1}
		}

		rate := c.HitCnt/c.ShowCnt + math.Sqrt((2*math.Log(totalShowCnt))/c.ShowCnt)
		if rate >= maxRate {
			maxRate = rate
			bestBannerID = bannerID
		}
	}

	return bestBannerID
}
