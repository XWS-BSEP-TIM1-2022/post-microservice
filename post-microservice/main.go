package main

import (
	"post-microservice/startup"
	cfg "post-microservice/startup/config"
)

func main() {
	config := cfg.NewConfig()
	server := startup.NewServer(config)
	server.Start()
}
