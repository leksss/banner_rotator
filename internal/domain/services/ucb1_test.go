package services

import (
	"testing"

	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/stretchr/testify/require"
)

const IterationCount = 10000

func TestUCB1BlankCounters(t *testing.T) {
	bannerIDs := []entities.BannerID{1, 2, 3, 4, 5}
	counters := entities.BannerCounterMap{}

	showStat := make(map[entities.BannerID]int)
	for i := 0; i < IterationCount; i++ {
		bannerID := CalculateBestBanner(bannerIDs, counters)
		require.True(t, bannerID != 0)

		counter := counters[bannerID]
		counter.ShowCnt++
		counters[bannerID] = counter

		showStat[bannerID]++
	}

	for _, bannerID := range bannerIDs {
		require.True(t, showStat[bannerID] > 1)
	}
}

func TestUCB1ZeroCounters(t *testing.T) {
	bannerIDs := []entities.BannerID{1, 2, 3, 4, 5}
	counters := entities.BannerCounterMap{
		1: {HitCnt: 0, ShowCnt: 0},
		2: {HitCnt: 0, ShowCnt: 0},
		3: {HitCnt: 0, ShowCnt: 0},
		4: {HitCnt: 0, ShowCnt: 0},
		5: {HitCnt: 0, ShowCnt: 0},
	}

	showStat := make(map[entities.BannerID]int)
	for i := 0; i < IterationCount; i++ {
		bannerID := CalculateBestBanner(bannerIDs, counters)
		require.True(t, bannerID != 0)

		counter := counters[bannerID]
		counter.ShowCnt++
		counters[bannerID] = counter

		showStat[bannerID]++
	}
	for _, bannerID := range bannerIDs {
		require.Equal(t, 2000, showStat[bannerID])
	}
}

func TestUCB1PopularBanner(t *testing.T) {
	bannerIDs := []entities.BannerID{1, 2, 3, 4, 5}
	counters := entities.BannerCounterMap{
		1: {HitCnt: 30, ShowCnt: 100},
		2: {HitCnt: 0, ShowCnt: 100},
		3: {HitCnt: 50, ShowCnt: 100}, // expected
		4: {HitCnt: 20, ShowCnt: 100},
		5: {HitCnt: 10, ShowCnt: 100},
	}

	showStat := make(map[entities.BannerID]int)
	for i := 0; i < IterationCount; i++ {
		bannerID := CalculateBestBanner(bannerIDs, counters)
		require.True(t, bannerID != 0)

		counter := counters[bannerID]
		counter.ShowCnt++
		counters[bannerID] = counter

		showStat[bannerID]++
	}

	require.Equal(t, entities.BannerID(3), maxRateBanner(showStat))
}

func maxRateBanner(showStat map[entities.BannerID]int) entities.BannerID {
	maxRate := 0
	bestBannerID := entities.BannerID(0)
	for bannerID, rate := range showStat {
		if rate >= maxRate {
			maxRate = rate
			bestBannerID = bannerID
		}
	}
	return bestBannerID
}
