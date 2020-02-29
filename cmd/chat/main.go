package main

import (
	"flag"
	"log"

	"github.com/MishaNiki/chat/internal/app/server"
)

var (
	pathConfig string
)

func init() {
	flag.StringVar(&pathConfig, "conf", "configs/server.json", "path to config file")
}

func main() {

	flag.Parse()

	config := server.NewConfig()
	if err := config.DecodeJSONConf(pathConfig); err != nil {
		log.Fatal(err)
	}
	s := server.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
