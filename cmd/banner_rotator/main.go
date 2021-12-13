package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/jmoiron/sqlx"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	"github.com/leksss/banner_rotator/internal/infrastructure/eventbus"
	"github.com/leksss/banner_rotator/internal/infrastructure/logger"
	mysql "github.com/leksss/banner_rotator/internal/infrastructure/storage/sql"
	grpc "github.com/leksss/banner_rotator/internal/server/grpc"
	"go.uber.org/zap"
)

//go:generate ../proto_generator.sh

func main() {
	configFile := flag.String("config", "configs/config.yaml", "path to conf file")
	flag.Parse()
	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	conf := config.NewConfig(*configFile)
	if err := conf.Parse(); err != nil {
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
	server.Start(ctx)
	<-ctx.Done()
	server.Stop(ctx)
}

func createKafkaConfig() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = 5
	return conf
}
