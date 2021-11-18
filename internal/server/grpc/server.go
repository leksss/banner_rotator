package internalgrpc

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/leksss/banner_rotator/internal/app"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
	pb "github.com/leksss/banner_rotator/proto/protobuf"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Server struct {
	grpcAddr string
	app      *app.App
	http     *http.Server
	grpc     *grpc.Server
	log      logger.Log
}

func NewServer(log logger.Log, app *app.App, config config.Config) *Server {
	return &Server{
		app:      app,
		log:      log,
		grpcAddr: fmt.Sprintf("%s:%s", config.GRPCAddr.Host, config.GRPCAddr.Port),
		http: &http.Server{
			Addr: fmt.Sprintf("%s:%s", config.HTTPAddr.Host, config.HTTPAddr.Port),
		},
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
	pb.RegisterBannerRotatorServiceServer(s.grpc, NewBannerRotatorService(s.app))

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

	s.http.Handler = logRequest(gwMux, s.log)
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

func logRequest(handler http.Handler, logger logger.Log) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		o := &responseObserver{ResponseWriter: w}
		handler.ServeHTTP(o, r)
		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}
		logger.Info(fmt.Sprintf("%s - - [%s] %q %d %d %q %q",
			addr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
			o.status,
			o.written,
			r.Referer(),
			r.UserAgent()))
	})
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	o.ResponseWriter.WriteHeader(code)
	if o.wroteHeader {
		return
	}
	o.wroteHeader = true
	o.status = code
}
