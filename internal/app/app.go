package app

import (
	"context"

	"github.com/leksss/banner-rotator/internal/domain/errors"
	"github.com/leksss/banner-rotator/internal/domain/interfaces"
	"github.com/leksss/banner-rotator/internal/infrastructure/logger"
	pb "github.com/leksss/banner-rotator/proto/protobuf"
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

	err := a.storage.HitBanner(ctx, in.SlotID, in.BannerID, in.GroupID)
	if err != nil {
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
		return &pb.GetBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotAndGroupAreRequired),
		}, nil
	}

	bannerID, err := a.storage.GetBanner(ctx, in.SlotID, in.GroupID)
	if err != nil {
		return &pb.GetBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.GetBannerResponse{
		Success:  true,
		BannerID: bannerID,
	}, nil
}

func toProtoError(errs []*pb.Error, err error) []*pb.Error {
	return append(errs, &pb.Error{
		Code: "banner",
		Msg:  err.Error(),
	})
}
