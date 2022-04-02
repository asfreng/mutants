package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"challenge/internal/handler"
	"challenge/internal/model"
	"challenge/internal/routers"
	"challenge/internal/util"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	config, err := util.LoadConfig("configs", "challenge")
	if err != nil {
		log.Printf("Could not read config file: %s", err.Error())
		return
	}

	// Conectar a Redis
	pool, err := util.ConnectToRedis(config.DatabaseConfig)
	if err != nil {
		fmt.Printf("Could not connect to Redis: %s\n", err.Error())
		return
	}
	defer util.CloseRedis(config.DatabaseConfig, pool)

	// Handler
	m := model.NewModel(pool, config)
	subHandler := handler.NewHandler(m)
	router := routers.SetupRouter(subHandler, config)

	// create a new server
	go router.Run(config.ServerConfig.BindAddress)

	// Wait forever, keeping main thread alive
	forever := make(chan bool)
	<-forever
}
