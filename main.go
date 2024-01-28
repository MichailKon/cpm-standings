package main

import (
	"cpm-standings/config"
	"cpm-standings/server_worker"
)

func main() {
	conf := config.ParseConfig("config.yaml")
	mapping := config.ParseStudentsHandlesMapping("mapping.yaml")
	server_worker.RunServer(conf, mapping)
}
