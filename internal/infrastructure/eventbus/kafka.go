package eventbus

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/leksss/banner_rotator/internal/domain/entities"
	"github.com/leksss/banner_rotator/internal/domain/interfaces"
)

type KafkaEventBus struct {
	conn  sarama.SyncProducer
	topic string
	log   interfaces.Log
}

func New(conn sarama.SyncProducer, topic string, log interfaces.Log) *KafkaEventBus {
	return &KafkaEventBus{
		conn:  conn,
		topic: topic,
		log:   log,
	}
}

func (k *KafkaEventBus) AddEvent(ctx context.Context, stat entities.EventStat) error {
	statJson, err := json.Marshal(stat)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(statJson),
	}
	_, _, err = k.conn.SendMessage(msg)
	if err != nil {
		return err
	}

	k.logEvent(msg)
	return nil
}

func (k *KafkaEventBus) logEvent(msg *sarama.ProducerMessage) {
	byteArg, _ := json.Marshal(msg)
	k.log.Info(fmt.Sprintf("kafka event: %s", string(byteArg)))
}
