package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcAddr string
	http     *http.Server
	grpc     *grpc.Server
	log      interfaces.Log
	storage  interfaces.Storage
	eventBus interfaces.EventBus
}

func NewServer(log interfaces.Log, config config.Config, storage interfaces.Storage,
	eventBus interfaces.EventBus) *Server {
	return &Server{
		log:      log,
		grpcAddr: config.GRPCAddr.DSN(),
		http: &http.Server{
			Addr: config.HTTPAddr.DSN(),
		},
		storage:  storage,
		eventBus: eventBus,
	}
}

func (s *Server) StartGRPC() error {
	lis, err := net.Listen("tcp", s.grpcAddr)
	if err != nil {
		s.log.Error("failed to listen:", zap.Error(err))
	}
	s.grpc = grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_zap.StreamServerInterceptor(s.log.GetLogger()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_zap.UnaryServerInterceptor(s.log.GetLogger()),
		)),
	)
	pb.RegisterBannerRotatorServiceServer(s.grpc, NewBannerRotatorService(s.log, s.storage, s.eventBus))

	s.log.Info(fmt.Sprintf("serving gRPC on %s", s.grpcAddr))
	return s.grpc.Serve(lis)
}

func (s *Server) StartHTTPProxy() error {
	conn, err := grpc.DialContext(
		context.Background(),
		s.grpcAddr,
		grpc.WithBlock(),
		grpc.WithInsecure(),
	)
	if err != nil {
		s.log.Error("failed to dial server:", zap.Error(err))
	}

	gwMux := runtime.NewServeMux()
	err = pb.RegisterBannerRotatorServiceHandler(context.Background(), gwMux, conn)
	if err != nil {
		s.log.Error("failed to register gateway:", zap.Error(err))
	}

	s.http.Handler = loggingMiddleware(gwMux, s.log)
	s.log.Info(fmt.Sprintf("serving gRPC-Gateway on %s", s.http.Addr))
	return s.http.ListenAndServe()
}

func (s *Server) StopHTTPProxy(ctx context.Context) error {
	s.log.Info("stopping HTTP proxy server...")
	return s.http.Shutdown(ctx)
}

func (s *Server) StopGRPC() {
	s.log.Info("stopping gRPC server...")
	s.grpc.GracefulStop()
}
