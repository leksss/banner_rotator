package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/leksss/banner_rotator/internal/domain/entities"
)

type KafkaConf struct {
	Host  string
	Port  string
	Topic string
}

type KafkaEventBus struct {
	conn   sarama.SyncProducer
	config KafkaConf
}

func New(config KafkaConf) *KafkaEventBus {
	return &KafkaEventBus{
		config: config,
	}
}

func (k *KafkaEventBus) Connect(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	dsn := fmt.Sprintf("%s:%s", k.config.Host, k.config.Port)
	conn, err := sarama.NewSyncProducer([]string{dsn}, config)
	if err != nil {
		return err
	}
	k.conn = conn
	return nil
}

func (k *KafkaEventBus) Close(ctx context.Context) error {
	return k.conn.Close()
}

func (k *KafkaEventBus) AddEvent(ctx context.Context, stat entities.EventStat) error {
	statJson, err := json.Marshal(stat)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: k.config.Topic,
		Value: sarama.StringEncoder(statJson),
	}
	_, _, err = k.conn.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
