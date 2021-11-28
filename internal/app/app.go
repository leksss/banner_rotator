package app

import (
	"context"
	"time"

	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/domain/errors"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/domain/services"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
)

type App struct {
	logger  logger.Log
	storage interfaces.Storage
	bus     interfaces.EventBus
}

func New(logger logger.Log, storage interfaces.Storage, bus interfaces.EventBus) *App {
	return &App{
		logger:  logger,
		storage: storage,
		bus:     bus,
	}
}

func (a *App) AddBanner(ctx context.Context, in *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	if in.SlotID == 0 || in.BannerID == 0 {
		return &pb.AddBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotAndBannerAreRequired),
		}, nil
	}

	if err := a.storage.AddBanner(ctx, in.SlotID, in.BannerID); err != nil {
		return &pb.AddBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.AddBannerResponse{
		Success: true,
	}, nil
}

func (a *App) RemoveBanner(ctx context.Context, in *pb.RemoveBannerRequest) (*pb.RemoveBannerResponse, error) {
	if in.SlotID == 0 || in.BannerID == 0 {
		return &pb.RemoveBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotAndBannerAreRequired),
		}, nil
	}

	if err := a.storage.RemoveBanner(ctx, in.SlotID, in.BannerID); err != nil {
		return &pb.RemoveBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.RemoveBannerResponse{
		Success: true,
	}, nil
}

func (a *App) HitBanner(ctx context.Context, in *pb.HitBannerRequest) (*pb.HitBannerResponse, error) {
	if in.SlotID == 0 || in.BannerID == 0 || in.GroupID == 0 {
		return &pb.HitBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotAndBannerAndGroupAreRequired),
		}, nil
	}

	incErrCh := a.goIncrementCounter(ctx, in.SlotID, in.BannerID, in.GroupID, true)
	if err := <-incErrCh; err != nil {
		return &pb.HitBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	busErrCh := a.goAddShowEvent(ctx, in.SlotID, in.BannerID, in.GroupID, entities.EventTypeHit)
	if err := <-busErrCh; err != nil {
		return &pb.HitBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.HitBannerResponse{
		Success: true,
	}, nil
}

func (a *App) GetBanner(ctx context.Context, in *pb.GetBannerRequest) (*pb.GetBannerResponse, error) {
	if in.SlotID == 0 || in.GroupID == 0 {
		return getBannerErrorResponse(errors.ErrInvalidRequestSlotAndGroupAreRequired)
	}

	bannerIDsCh, banErrCh := a.goGetBannersBySlot(ctx, in.SlotID)
	countersCh, cntErrCh := a.goGetSlotCounters(ctx, in.SlotID, in.GroupID)
	if err := <-banErrCh; err != nil {
		return getBannerErrorResponse(err)
	}
	if err := <-cntErrCh; err != nil {
		return getBannerErrorResponse(err)
	}

	bannerIDs := <-bannerIDsCh
	if len(bannerIDs) == 0 {
		return getBannerErrorResponse(errors.ErrNoAvailableBannersInSlot)
	}
	counters := <-countersCh

	bestBannerID := services.CalculateBestBanner(a.logger, bannerIDs, counters)
	if bestBannerID == 0 {
		return getBannerErrorResponse(errors.ErrBannerNotFound)
	}

	incErrCh := a.goIncrementCounter(ctx, in.SlotID, uint64(bestBannerID), in.GroupID, false)
	busErrCh := a.goAddShowEvent(ctx, in.SlotID, uint64(bestBannerID), in.GroupID, entities.EventTypeShow)
	if err := <-incErrCh; err != nil {
		return getBannerErrorResponse(err)
	}
	if err := <-busErrCh; err != nil {
		return getBannerErrorResponse(err)
	}

	return &pb.GetBannerResponse{
		Success:  true,
		BannerID: uint64(bestBannerID),
	}, nil
}

func toProtoError(errs []*pb.Error, err error) []*pb.Error {
	return append(errs, &pb.Error{
		Code: "banner",
		Msg:  err.Error(),
	})
}

func getBannerErrorResponse(err error) (*pb.GetBannerResponse, error) {
	return &pb.GetBannerResponse{
		Success: false,
		Errors:  toProtoError([]*pb.Error{}, err),
	}, nil
}

func (a *App) goGetBannersBySlot(ctx context.Context, slotID uint64) (<-chan []entities.BannerID, <-chan error) {
	resCh := make(chan []entities.BannerID, 1)
	errCh := make(chan error, 1)
	go func() {
		bannerIDs, err := a.storage.GetBannersBySlot(ctx, slotID)
		errCh <- err
		resCh <- bannerIDs
	}()
	return resCh, errCh
}

func (a *App) goGetSlotCounters(ctx context.Context, slotID, groupID uint64) (<-chan entities.BannerCounterMap, <-chan error) {
	resCh := make(chan entities.BannerCounterMap, 1)
	errCh := make(chan error, 1)
	go func() {
		counters, err := a.storage.GetSlotCounters(ctx, slotID, groupID)
		errCh <- err
		resCh <- counters
	}()
	return resCh, errCh
}

func (a *App) goIncrementCounter(ctx context.Context, slotID uint64, bannerID uint64, groupID uint64, isHit bool) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		var err error
		if isHit {
			err = a.storage.IncrementHit(ctx, slotID, bannerID, groupID)
		} else {
			err = a.storage.IncrementShow(ctx, slotID, bannerID, groupID)
		}
		errCh <- err
	}()
	return errCh
}

func (a *App) goAddShowEvent(ctx context.Context, slotID uint64, bannerID uint64, groupID uint64, eventType entities.EventType) <-chan error {
	errCh := make(chan error, 1)
	go func() {
		stat := entities.EventStat{
			EventType: eventType,
			SlotID:    slotID,
			BannerID:  bannerID,
			GroupID:   groupID,
			CreatedAt: time.Now(),
		}
		err := a.bus.AddEvent(ctx, stat)
		errCh <- err
	}()
	return errCh
}
