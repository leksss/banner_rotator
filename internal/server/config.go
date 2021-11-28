package server

import "fmt"

type Config struct {
	Host string
	Port string
}

func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
