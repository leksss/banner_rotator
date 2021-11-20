package services

import (
	"fmt"
	"math"

	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
)

// CalculateBestBanner calculate Upper Confidence Bound algorithm
func CalculateBestBanner(log logger.Log, counters []*entities.Counter) uint64 {
	totalShowCnt := float64(0)
	for _, c := range counters {
		totalShowCnt += c.ShowCnt
	}

	rates := make(map[uint64]float64)
	for _, c := range counters {
		rates[c.BannerID] = c.HitCnt/c.ShowCnt + math.Sqrt((2*math.Log(totalShowCnt))/c.ShowCnt)
		log.Info(fmt.Sprintf("rate for bannerID %d: %f", c.BannerID, rates[c.BannerID]))
	}

	bestBannerID := uint64(0)
	maxRate := float64(0)
	for bannerID, rate := range rates {
		if maxRate < rate {
			bestBannerID = bannerID
			maxRate = rate
		}
	}

	return bestBannerID
}
