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

func TestChessSuite(t *testing.T) {
	suite.Run(t, new(BannerRotatorSuite))
}

func (s *BannerRotatorSuite) TestAddTheSameBannerTwice() {
	request := &pb.AddBannerRequest{
		SlotID:   1,
		BannerID: 1,
	}
	response, _ := s.client.AddBanner(s.ctx, request)
	s.Require().True(response.Success)
	s.Require().Equal(0, len(response.Errors))

	response, _ = s.client.AddBanner(s.ctx, request)
	s.Require().False(response.Success)
	s.Require().Equal(1, len(response.Errors))
}

func (s *BannerRotatorSuite) TestRemoveBanner() {
	addResponse, _ := s.client.AddBanner(s.ctx, &pb.AddBannerRequest{
		SlotID:   1,
		BannerID: 1,
	})
	s.Require().True(addResponse.Success)

	getResponse, _ := s.client.GetBanner(s.ctx, &pb.GetBannerRequest{
		SlotID:  1,
		GroupID: 1,
	})
	s.Require().True(getResponse.Success)
	s.Require().Equal(uint64(1), getResponse.BannerID)

	removeResponse, _ := s.client.RemoveBanner(s.ctx, &pb.RemoveBannerRequest{
		SlotID:   1,
		BannerID: 1,
	})
	s.Require().True(getResponse.Success)
	s.Require().Equal(0, len(removeResponse.Errors))

	getResponse, _ = s.client.GetBanner(s.ctx, &pb.GetBannerRequest{
		SlotID:  1,
		GroupID: 1,
	})
	s.Require().False(getResponse.Success)
	s.Require().Equal(uint64(0), getResponse.BannerID)
}

//func (s *BannerRotatorSuite) TestHitBanner() {
//
//}
