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

	"github.com/Shopify/sarama"
	"github.com/jmoiron/sqlx"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	"github.com/leksss/banner_rotator/internal/infrastructure/eventbus"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
	mysql "github.com/leksss/banner_rotator/internal/infrastructure/storage/sql"
	grpc "github.com/leksss/banner_rotator/internal/server/grpc"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

//go:generate ../proto_generator.sh

const (
	appShutdownMessage      = "application exits"
	gracefulShutdownTimeout = 3 * time.Second
)

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to conf file")
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf := config.NewConfig(*configFile)
	err := conf.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolintlint
	}

	var zapConfig zap.Config
	if conf.IsDebug() {
		zapConfig = zap.NewDevelopmentConfig()
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	logg := logger.New(zapConfig, conf.Logger, conf.GetProjectRoot())

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	dbConn, err := sqlx.ConnectContext(ctx, "mysql", conf.Database.DSN())
	if err != nil {
		logg.Error(fmt.Sprintf("connect to storage failed: %s", err.Error()))
	}
	defer dbConn.Close()

	kafkaConn, err := sarama.NewSyncProducer([]string{conf.Kafka.DSN()}, createKafkaConfig())
	if err != nil {
		logg.Error(fmt.Sprintf("connect to event bus failed: %s", err.Error()))
	}
	defer kafkaConn.Close()

	storage := mysql.New(dbConn, logg)
	bus := eventbus.New(kafkaConn, conf.Kafka.Topic, logg)

	server := grpc.NewServer(logg, conf, storage, bus)

	errs := make(chan error)

	go func() {
		errs <- server.StartGRPC()
	}()

	go func() {
		time.Sleep(500 * time.Millisecond)
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

func createKafkaConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 5
	return conf
}
