package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leksss/banner_rotator/internal/infrastructure/config"
	"github.com/pressly/goose/v3"
)

/*
@Usage
$ go run . -dir=./migrations mysql status
$ go run . -dir=./migrations mysql create init sql
$ go run . -dir=./migrations mysql up
*/

var (
	flags      = flag.NewFlagSet("goose", flag.ExitOnError)
	dir        = flags.String("dir", ".", "directory with migration files")
	configFile = flags.String("conf", "configs/config.yaml", "path to conf file")
)

func main() {
	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()
		return
	}

	fmt.Println(args)

	conf := config.NewConfig(*configFile)
	err := conf.Parse()
	if err != nil {
		log.Fatal(err.Error()) //nolintlint
	}
	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true",
		conf.Database.User, conf.Database.Password, conf.Database.Name)

	db, err := goose.OpenDBWithDriver("mysql", dsn)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	command := args[1] //nolint
	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err) //nolint
	}
}
