package sqlstorage

import "fmt"

type DatabaseConf struct {
	Host     string
	User     string
	Password string
	Name     string
}

func (c *DatabaseConf) DSN() string {
	return fmt.Sprintf("%s:%s@(%s:3306)/%s?parseTime=true", c.User, c.Password, c.Host, c.Name)
}
