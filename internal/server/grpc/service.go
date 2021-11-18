package internalgrpc

import (
	"context"

	"github.com/leksss/banner_rotator/internal/app"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
)

type BannerRotatorService struct {
	pb.UnimplementedBannerRotatorServiceServer
	app *app.App
}

func NewBannerRotatorService(app *app.App) *BannerRotatorService {
	return &BannerRotatorService{app: app}
}

func (s *BannerRotatorService) AddBanner(ctx context.Context, in *pb.AddBannerRequest) (*pb.AddBannerResponse, error) {
	return s.app.AddBanner(ctx, in)
}

func (s *BannerRotatorService) RemoveBanner(ctx context.Context, in *pb.RemoveBannerRequest) (*pb.RemoveBannerResponse, error) {
	return s.app.RemoveBanner(ctx, in)
}

func (s *BannerRotatorService) HitBanner(ctx context.Context, in *pb.HitBannerRequest) (*pb.HitBannerResponse, error) {
	return s.app.HitBanner(ctx, in)
}

func (s *BannerRotatorService) GetBanner(ctx context.Context, in *pb.GetBannerRequest) (*pb.GetBannerResponse, error) {
	return s.app.GetBanner(ctx, in)
}
