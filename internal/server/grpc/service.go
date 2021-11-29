package internalgrpc

import (
	"context"
	"time"

	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/domain/errors"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/domain/services"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
)

type BannerRotatorService struct {
	pb.UnimplementedBannerRotatorServiceServer
	log      interfaces.Log
	storage  interfaces.Storage
	eventBus interfaces.EventBus
}

func NewBannerRotatorService(log interfaces.Log, storage interfaces.Storage,
	eventBus interfaces.EventBus) *BannerRotatorService {
	return &BannerRotatorService{
		log:      log,
		storage:  storage,
		eventBus: eventBus,
	}
}

func (s *BannerRotatorService) AddBanner(ctx context.Context, in *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	if in.SlotID == 0 || in.BannerID == 0 {
		return &pb.AddBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotBannerRequired),
		}, nil
	}

	if err := s.storage.AddBanner(ctx, in.SlotID, in.BannerID); err != nil {
		return &pb.AddBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.AddBannerResponse{
		Success: true,
	}, nil
}

func (s *BannerRotatorService) RemoveBanner(ctx context.Context,
	in *pb.RemoveBannerRequest) (*pb.RemoveBannerResponse, error) {
	if in.SlotID == 0 || in.BannerID == 0 {
		return &pb.RemoveBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotBannerRequired),
		}, nil
	}

	if err := s.storage.RemoveBanner(ctx, in.SlotID, in.BannerID); err != nil {
		return &pb.RemoveBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	return &pb.RemoveBannerResponse{
		Success: true,
	}, nil
}

func (s *BannerRotatorService) HitBanner(ctx context.Context, in *pb.HitBannerRequest) (*pb.HitBannerResponse, error) {
	if in.SlotID == 0 || in.BannerID == 0 || in.GroupID == 0 {
		return &pb.HitBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, errors.ErrInvalidRequestSlotBannerGroupRequired),
		}, nil
	}

	err := s.storage.IncrementHit(ctx, in.SlotID, in.BannerID, in.GroupID)
	if err != nil {
		return &pb.HitBannerResponse{
			Success: false,
			Errors:  toProtoError([]*pb.Error{}, err),
		}, nil
	}

	err = s.eventBus.AddEvent(ctx, entities.EventStat{
		EventType: entities.EventTypeHit,
		SlotID:    in.SlotID,
		BannerID:  in.BannerID,
		GroupID:   in.GroupID,
		CreatedAt: time.Now(),
	})
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

func (s *BannerRotatorService) GetBanner(ctx context.Context, in *pb.GetBannerRequest) (*pb.GetBannerResponse, error) {
	if in.SlotID == 0 || in.GroupID == 0 {
		return getBannerErrorResponse(errors.ErrInvalidRequestSlotGroupRequired)
	}

	bannerIDs, err := s.storage.GetBannersBySlot(ctx, in.SlotID)
	if err != nil {
		return getBannerErrorResponse(err)
	}
	if len(bannerIDs) == 0 {
		return getBannerErrorResponse(errors.ErrNoAvailableBannersInSlot)
	}

	counters, err := s.storage.GetSlotCounters(ctx, in.SlotID, in.GroupID)
	if err != nil {
		return getBannerErrorResponse(err)
	}

	bestBannerID := services.CalculateBestBanner(bannerIDs, counters)
	if bestBannerID == 0 {
		return getBannerErrorResponse(errors.ErrBannerNotFound)
	}
	bannerID := uint64(bestBannerID)

	err = s.storage.IncrementShow(ctx, in.SlotID, bannerID, in.GroupID)
	if err != nil {
		return getBannerErrorResponse(err)
	}

	err = s.eventBus.AddEvent(ctx, entities.EventStat{
		EventType: entities.EventTypeShow,
		SlotID:    in.SlotID,
		BannerID:  bannerID,
		GroupID:   in.GroupID,
		CreatedAt: time.Now(),
	})
	if err != nil {
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
