// +build integration

package integration_test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
)

type BannerRotatorSuite struct {
	suite.Suite
	ctx    context.Context
	config config.Config
	db     *sqlx.DB
	client pb.BannerRotatorServiceClient
}

func (s *BannerRotatorSuite) SetupSuite() {
	s.config = config.NewConfig("configs/config.yaml")
	err := s.config.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolintlint
	}

	grpcConn, err := grpc.Dial(s.config.GRPCAddr.DSN(), grpc.WithInsecure())
	s.Require().NoError(err)

	s.ctx = context.Background()
	s.client = pb.NewBannerRotatorServiceClient(grpcConn)
}

func (s *BannerRotatorSuite) SetupTest() {
	var seed int64 = time.Now().UnixNano()
	rand.Seed(seed)
	s.T().Log("seed:", seed)

	var err error
	s.db, err = sqlx.ConnectContext(s.ctx, "mysql", s.config.Database.DSN())
	if err != nil {
		log.Fatal(fmt.Sprintf("connect to storage failed: %s", err.Error()))
	}

	s.db.NamedExecContext(s.ctx, "TRUNCATE slot2banner", make(map[string]interface{}))
	s.db.NamedExecContext(s.ctx, "TRUNCATE ucb1", make(map[string]interface{}))
}

func (s *BannerRotatorSuite) TearDownTest() {
	s.db.Close()
}

func TestBannerRotatorSuite(t *testing.T) {
	suite.Run(t, new(BannerRotatorSuite))
}

func (s *BannerRotatorSuite) TestAddTheSameBannerTwice() {
	request := &pb.AddBannerRequest{
		SlotID:   1,
		BannerID: 1,
	}
	response, _ := s.client.AddBanner(s.ctx, request)
	s.Require().True(response.GetSuccess())
	s.Require().Equal(0, len(response.GetErrors()))

	response, _ = s.client.AddBanner(s.ctx, request)
	s.Require().False(response.GetSuccess())
	s.Require().Equal(1, len(response.GetErrors()))
}

func (s *BannerRotatorSuite) TestRemoveBanner() {
	addResponse, _ := s.client.AddBanner(s.ctx, &pb.AddBannerRequest{
		SlotID:   1,
		BannerID: 1,
	})
	s.Require().True(addResponse.GetSuccess())

	getResponse, _ := s.client.GetBanner(s.ctx, &pb.GetBannerRequest{
		SlotID:  1,
		GroupID: 1,
	})
	s.Require().True(getResponse.GetSuccess())
	s.Require().Equal(uint64(1), getResponse.GetBannerID())

	removeResponse, _ := s.client.RemoveBanner(s.ctx, &pb.RemoveBannerRequest{
		SlotID:   1,
		BannerID: 1,
	})
	s.Require().True(getResponse.GetSuccess())
	s.Require().Equal(0, len(removeResponse.GetErrors()))

	getResponse, _ = s.client.GetBanner(s.ctx, &pb.GetBannerRequest{
		SlotID:  1,
		GroupID: 1,
	})
	s.Require().False(getResponse.GetSuccess())
	s.Require().Equal(uint64(0), getResponse.GetBannerID())
}

func (s *BannerRotatorSuite) TestGetAndHitBanner() {
	bannersCount := uint64(3)

	var i uint64
	for i = 1; i <= bannersCount; i++ {
		addResponse, _ := s.client.AddBanner(s.ctx, &pb.AddBannerRequest{
			SlotID:   1,
			BannerID: i,
		})
		s.Require().True(addResponse.GetSuccess())
	}

	bannerMap := make(map[uint64]uint64)
	for i = 1; i <= 10; i++ {
		getResponse, _ := s.client.GetBanner(s.ctx, &pb.GetBannerRequest{
			SlotID:  1,
			GroupID: 1,
		})
		s.Require().True(getResponse.GetSuccess())
		bannerMap[getResponse.GetBannerID()] = getResponse.GetBannerID()
	}
	s.Require().Equal(int(bannersCount), len(bannerMap))

	favoriteBannerID := uint64(2)
	for i = 1; i <= 10; i++ {
		getResponse, _ := s.client.HitBanner(s.ctx, &pb.HitBannerRequest{
			SlotID:   1,
			GroupID:  1,
			BannerID: favoriteBannerID,
		})
		s.Require().True(getResponse.GetSuccess())
	}

	for i = 1; i <= 10; i++ {
		getResponse, _ := s.client.GetBanner(s.ctx, &pb.GetBannerRequest{
			SlotID:  1,
			GroupID: 1,
		})
		s.Require().True(getResponse.GetSuccess())
		s.Require().Equal(favoriteBannerID, getResponse.GetBannerID())
	}
}
