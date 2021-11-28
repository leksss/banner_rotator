package eventbus

import "fmt"

type KafkaConf struct {
	Host  string
	Port  string
	Topic string
}

func (c *KafkaConf) DSN() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
