package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	restartTimeout = time.Second
)

type Server struct {
	grpcAddr string
	wg       *sync.WaitGroup
	http     *http.Server
	grpc     *grpc.Server
	log      interfaces.Log
	storage  interfaces.Storage
	eventBus interfaces.EventBus
}

func NewServer(log interfaces.Log,
	config config.Config, storage interfaces.Storage, eventBus interfaces.EventBus) interfaces.StartStopper {
	return &Server{
		log:      log,
		grpcAddr: config.GRPCAddr.DSN(),
		http: &http.Server{
			Addr: config.HTTPAddr.DSN(),
		},
		storage:  storage,
		eventBus: eventBus,
		wg:       &sync.WaitGroup{},
	}
}

func (s *Server) Start(ctx context.Context) {
	s.wg.Add(1)
	go func(ctx context.Context) {
		defer s.wg.Done()
		s.startGRPC(ctx)
	}(ctx)

	s.wg.Add(1)
	go func(ctx context.Context) {
		defer s.wg.Done()
		s.startHTTPProxy(ctx)
	}(ctx)
}

func (s *Server) Stop(ctx context.Context) {
	s.stopHTTPProxy(ctx)
	s.stopGRPC(ctx)
	s.wg.Wait()
}

func (s *Server) startGRPC(ctx context.Context) {
	var lis net.Listener
	var err error
	for {
		if lis, err = net.Listen("tcp", s.grpcAddr); err == nil {
			break
		}
		s.log.Warn(fmt.Sprint("gRPC listener failed — restarting in ", restartTimeout), zap.Error(err))
		select {
		case <-ctx.Done():
			break
		case <-time.After(restartTimeout):
			continue
		}
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
	if err := s.grpc.Serve(lis); err != nil {
		s.log.Error("serving gRPC failed", zap.Error(err))
	}
}

func (s *Server) startHTTPProxy(ctx context.Context) {
	var conn *grpc.ClientConn
	var err error
	for {
		if conn, err = grpc.DialContext(ctx, s.grpcAddr, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
			break
		}
		s.log.Warn(fmt.Sprint("failed to dial gRPC server — restarting in ", restartTimeout), zap.Error(err))
		select {
		case <-ctx.Done():
			break
		case <-time.After(restartTimeout):
			continue
		}
	}

	gwMux := runtime.NewServeMux()
	err = pb.RegisterBannerRotatorServiceHandler(context.Background(), gwMux, conn)
	if err != nil {
		s.log.Error("failed to register gateway:", zap.Error(err))
	}

	s.http.Handler = loggingMiddleware(gwMux, s.log)

	s.log.Info(fmt.Sprintf("serving gRPC-Gateway on %s", s.http.Addr))
	if err := s.http.ListenAndServe(); err != nil {
		if !errors.Is(ctx.Err(), context.Canceled) {
			s.log.Error("serving gRPC-Gateway failed", zap.Error(err))
		}
	}
}

func (s *Server) stopHTTPProxy(ctx context.Context) {
	s.log.Info("stopping gRPC-Gateway...")
	if err := s.http.Shutdown(ctx); err != nil {
		if errors.Is(err, context.Canceled) {
			s.log.Info("gRPC-Gateway has been stopped by canceled context")
		} else {
			s.log.Error("stopping gRPC-Gateway failed", zap.Error(err))
		}
	}
}

func (s *Server) stopGRPC(ctx context.Context) {
	s.log.Info("stopping gRPC server...")
	s.grpc.GracefulStop()
	if errors.Is(ctx.Err(), context.Canceled) {
		s.log.Info("gRPC server has been stopped by canceled context")
	}
}
