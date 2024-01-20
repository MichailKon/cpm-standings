package main

import (
	"cpm-standings/config"
	"cpm-standings/server_worker"
)

func main() {
	conf := config.ParseConfig("config.yaml")
	server_worker.RunServer(conf)
}
