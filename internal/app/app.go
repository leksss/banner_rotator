package app

import (
	"context"

	"github.com/leksss/banner_rotator/internal/domain/errors"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/domain/services"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
)

type App struct {
	logger  logger.Log
	storage interfaces.Storage
}

func New(logger logger.Log, storage interfaces.Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
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

	if err := a.storage.IncrementHit(ctx, in.SlotID, in.BannerID, in.GroupID); err != nil {
		return &pb.HitBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	//TODO отправляем событие хита в очередь для аналитической системы

	return &pb.HitBannerResponse{
		Success: true,
	}, nil
}

func (a *App) GetBanner(ctx context.Context, in *pb.GetBannerRequest) (*pb.GetBannerResponse, error) {
	if in.SlotID == 0 || in.GroupID == 0 {
		return &pb.GetBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotAndGroupAreRequired),
		}, nil
	}

	bannerIDs, err := a.storage.GetBannersBySlot(ctx, in.SlotID)
	if err != nil {
		return nil, err
	}
	if len(bannerIDs) == 0 {
		return &pb.GetBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrNoAvailableBannersInSlot),
		}, nil
	}

	counters, err := a.storage.GetSlotCounters(ctx, in.SlotID, in.GroupID)
	if err != nil {
		return nil, err
	}

	bestBannerID := services.CalculateBestBanner(a.logger, bannerIDs, counters)
	if bestBannerID == 0 {
		return &pb.GetBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrBannerNotFound),
		}, nil
	}

	if err := a.storage.IncrementShow(ctx, in.SlotID, bestBannerID, in.GroupID); err != nil {
		return &pb.GetBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	// TODO отправляем событие показа в очередь для аналитической системы

	return &pb.GetBannerResponse{
		Success:  true,
		BannerID: bestBannerID,
	}, nil
}

func toProtoError(errs []*pb.Error, err error) []*pb.Error {
	return append(errs, &pb.Error{
		Code: "banner",
		Msg:  err.Error(),
	})
}
