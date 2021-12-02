package sqlstorage

import "fmt"

type DatabaseConf struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

func (c *DatabaseConf) DSN() string {
	return fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true", c.User, c.Password, c.Host, c.Port, c.Name)
}
