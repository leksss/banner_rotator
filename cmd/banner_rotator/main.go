package main

import (
	"context"
	"flag"
	"fmt"
	"log"
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
	appShutdownMessage = "application exits"
	serverDownTimeout  = 1 * time.Second
	serverStartTimeout = 500 * time.Millisecond
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
	serviceStart(ctx, server, errs)
	serviceStop(ctx, server, errs)

	for err := range errs {
		if err == nil {
			continue
		}
		logg.Info("shutdown err message", zap.Error(err))
		if err.Error() == appShutdownMessage {
			return
		}
	}
}

func serviceStart(ctx context.Context, server *grpc.Server, errs chan<- error) {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), serverStartTimeout)
	defer cancel()

	go func() {
		errs <- server.StartGRPC()
	}()

	<-timeoutCtx.Done()

	go func(ctx context.Context) {
		errs <- server.StartHTTPProxy(ctx)
	}(ctx)
}

func serviceStop(ctx context.Context, server *grpc.Server, errs chan<- error) {
	<-ctx.Done()

	timeoutCtx, cancel := context.WithTimeout(context.Background(), serverDownTimeout)
	defer cancel()

	go func(ctx context.Context) {
		errs <- server.StopHTTPProxy(ctx)
	}(timeoutCtx)

	<-timeoutCtx.Done()

	go func() {
		server.StopGRPC()
		errs <- errors.New(appShutdownMessage)
	}()
}

func createKafkaConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 5
	return conf
}
