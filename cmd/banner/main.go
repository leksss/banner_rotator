package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/leksss/banner_rotator/internal/app"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
	memory "github.com/leksss/banner_rotator/internal/infrastructure/storage/memory"
	mysql "github.com/leksss/banner_rotator/internal/infrastructure/storage/sql"
	grpc "github.com/leksss/banner_rotator/internal/server/grpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate ./proto_generator.sh

const (
	appShutdownMessage      = "application exits"
	gracefulShutdownTimeout = 3 * time.Second
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to conf file")
	conf := config.NewConfig(*configFile)
	err := conf.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolintlint
	}

	flag.Parse()
	//if flag.Arg(0) == "version" {
	//	printVersion()
	//	return
	//}

	logg := logger.New(conf.Logger, conf.GetProjectRoot(), conf.IsDebug())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	storage := createStorageInstance(ctx, conf, logg)
	defer storage.Close(ctx)

	calendar := app.New(logg, storage)
	server := grpc.NewServer(logg, calendar, conf)

	errs := make(chan error)

	go func() {
		errs <- server.StartGRPC()
	}()

	go func() {
		errs <- server.StartHTTPProxy()
	}()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)
		sig := <-quit

		logg.Warn("os signal received, beginning graceful shutdown with timeout",
			zap.String("signal", sig.String()),
			zap.Duration("timeout", gracefulShutdownTimeout),
		)

		success := make(chan string)
		go func() {
			errs <- server.StopHTTPProxy(context.Background())
			success <- "HTTP server successfully stopped"
		}()
		go func() {
			server.StopGRPC()
			success <- "gRPC server successfully stopped"
		}()
		go func() {
			time.Sleep(gracefulShutdownTimeout)
			logg.Error("failed to gracefully shut down server within timeout. Shutting down with Fatal",
				zap.Duration("timeout", gracefulShutdownTimeout))
		}()
		logg.Info(<-success)
		logg.Info(<-success)
		errs <- errors.New(appShutdownMessage)
	}()

	for err := range errs {
		if err == nil {
			continue
		}
		logg.Warn("shutdown err message", zap.Error(err))
		if err.Error() == appShutdownMessage {
			return
		}
	}
}

func createStorageInstance(ctx context.Context, conf config.Config, logg logger.Log) interfaces.Storage {
	var storage interfaces.Storage
	if conf.Env == config.EnvTest {
		storage = memory.New()
	} else {
		storage = mysql.New(conf.Database, logg)
	}

	if err := storage.Connect(ctx); err != nil {
		logg.Error(fmt.Sprintf("Connect to storage failed: %s", err.Error()))
	}
	return storage
}
